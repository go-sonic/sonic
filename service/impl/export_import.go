package impl

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cast"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v2"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/pageparser"
	"github.com/go-sonic/sonic/util/xerr"
)

type exportImport struct {
	CategoryService     service.CategoryService
	PostService         service.PostService
	TagService          service.TagService
	PostTagService      service.PostTagService
	PostCategoryService service.PostCategoryService
}

func NewExportImport(categoryService service.CategoryService,
	postService service.PostService,
	tagService service.TagService,
	postTagService service.PostTagService,
	postCategoryService service.PostCategoryService,
) service.ExportImport {
	return &exportImport{
		CategoryService:     categoryService,
		PostService:         postService,
		TagService:          tagService,
		PostTagService:      postTagService,
		PostCategoryService: postCategoryService,
	}
}

func (e *exportImport) CreateByMarkdown(ctx context.Context, filename string, reader io.Reader) (*entity.Post, error) {
	contentFrontMatter, err := pageparser.ParseFrontMatterAndContent(reader)
	if err != nil {
		return nil, xerr.WithMsg(err, "parse markdown failed").WithStatus(xerr.StatusInternalServerError)
	}

	content, frontmatter := string(contentFrontMatter.Content), contentFrontMatter.FrontMatter

	postDate, postName, err := parseJekyllFilename(filename)
	if err == nil {
		content = convertJekyllContent(frontmatter, content)
		frontmatter = convertJekyllMetaData(frontmatter, postName, postDate)
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(content), &buf); err != nil {
		return nil, xerr.BadParam.Wrapf(err, "convert markdown err").WithStatus(xerr.StatusBadRequest)
	}

	post := param.Post{
		Status:          consts.PostStatusPublished,
		EditorType:      consts.EditorTypeMarkdown.Ptr(),
		OriginalContent: content,
		Content:         buf.String(),
	}

	for key, value := range frontmatter {
		switch key {
		case "title":
			post.Title = value.(string)
		case "permalink":
			post.Slug = value.(string)
		case "slug":
			post.Slug = value.(string)
		case "date":
			if s, ok := value.(string); ok {
				date, err := cast.StringToDate(s)
				if err != nil {
					log.CtxWarnf(ctx, "CreateByMarkdown convert date time err=%v", err)
				} else {
					post.CreateTime = util.Int64Ptr(date.UnixMilli())
				}
			}
		case "summary":
			post.Summary = value.(string)
		case "draft":
			if s, ok := value.(string); ok && s == "true" {
				post.Status = consts.PostStatusDraft
			}
			if b, ok := value.(bool); ok && b {
				post.Status = consts.PostStatusDraft
			}
		case "updated":
			if s, ok := value.(string); ok {
				date, err := cast.StringToDate(s)
				if err != nil {
					log.CtxWarnf(ctx, "CreateByMarkdown convert updated time err=%v", err)
				} else {
					post.UpdateTime = util.Int64Ptr(date.UnixMilli())
				}
			}
		case "lastmod":
			if s, ok := value.(string); ok {
				date, err := cast.StringToDate(s)
				if err != nil {
					log.CtxWarnf(ctx, "CreateByMarkdown convert lastmod time err=%v", err)
				} else {
					post.EditTime = util.Int64Ptr(date.UnixMilli())
				}
			}
		case "keywords":
			post.MetaKeywords = value.(string)
		case "comments":
			if s, ok := value.(string); ok && s == "true" {
				post.DisallowComment = false
			}
			if b, ok := value.(bool); ok && b {
				post.DisallowComment = !b
			}
			switch s := value.(type) {
			case string:
				comments, err := strconv.ParseBool(s)
				if err != nil {
					log.CtxWarnf(ctx, "CreateByMarkdown parse comments err=%v", err)
				} else {
					post.DisallowComment = !comments
				}
			case bool:
				post.DisallowComment = !s
			}
		case "tags":
			if _, ok := value.([]any); !ok {
				continue
			}

			for _, v := range value.([]any) {
				if _, ok := v.(string); !ok {
					continue
				}
				tag, err := e.TagService.GetByName(ctx, v.(string))
				if err != nil && xerr.GetType(err) == xerr.NoRecord {
					tag, err := e.TagService.Create(ctx, &param.Tag{
						Name: v.(string),
						Slug: util.Slug(v.(string)),
					})
					if err != nil {
						post.TagIDs = append(post.TagIDs, tag.ID)
					}
				} else if err == nil {
					post.TagIDs = append(post.TagIDs, tag.ID)
				}
			}

		case "categories", "category":

			switch s := value.(type) {
			case string:
				// example:
				// ---
				// categories: life
				// ---
				name := strings.TrimSpace(s)
				category, err := e.CategoryService.GetByName(ctx, name)
				switch {
				case xerr.GetType(err) == xerr.NoRecord:
					categoryParam := &param.Category{
						Name: name,
						Slug: util.Slug(name),
					}
					category, err = e.CategoryService.Create(ctx, categoryParam)
					if err != nil {
						return nil, err
					}
					post.CategoryIDs = append(post.CategoryIDs, category.ID)
				case err != nil:
					return nil, err
				case err == nil:
					post.CategoryIDs = append(post.CategoryIDs, category.ID)
				}
			case []any:
				// example:
				// ---
				// categories:
				// - Development
				// - VIM
				// ---
				// VIM is sub category of Development
				var parentCategoryID int32
				for _, v := range s {
					if _, ok := v.(string); !ok {
						continue
					}
					category, err := e.CategoryService.GetByName(ctx, v.(string))
					switch {
					case xerr.GetType(err) == xerr.NoRecord:
						categoryParam := &param.Category{
							Name:     v.(string),
							Slug:     util.Slug(v.(string)),
							ParentID: parentCategoryID,
						}
						category, err = e.CategoryService.Create(ctx, categoryParam)
						if err != nil {
							return nil, err
						}
						post.CategoryIDs = append(post.CategoryIDs, category.ID)
						parentCategoryID = category.ID
					case err != nil:
						return nil, err
					case err == nil:
						post.CategoryIDs = append(post.CategoryIDs, category.ID)
					}
				}
			}
		}
	}
	return e.PostService.Create(ctx, &post)
}

func (e *exportImport) ExportMarkdown(ctx context.Context, needFrontMatter bool) (string, error) {
	posts, _, err := e.PostService.Page(ctx, param.PostQuery{
		Page: param.Page{
			PageNum:  0,
			PageSize: 999999,
		},
		Statuses: []*consts.PostStatus{consts.PostStatusDraft.Ptr(), consts.PostStatusIntimate.Ptr(), consts.PostStatusPublished.Ptr()},
	})
	if err != nil {
		return "", err
	}

	backupFilename := consts.SonicBackupMarkdownPrefix + time.Now().Format("2006-01-02-15-04-05") + util.GenUUIDWithOutDash() + ".zip"
	backupFilePath := config.BackupMarkdownDir

	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		err = os.MkdirAll(backupFilePath, os.ModePerm)
		if err != nil {
			return "", xerr.NoType.Wrap(err).WithMsg("create dir err")
		}
	} else if err != nil {
		return "", xerr.NoType.Wrap(err).WithMsg("get fileInfo")
	}

	toBackupPaths := []string{}
	for _, post := range posts {
		var markdown strings.Builder
		if needFrontMatter {
			frontMatter, err := e.getFrontMatterYaml(ctx, post)
			if err == nil {
				markdown.WriteString("---\n")
				markdown.WriteString(frontMatter)
				markdown.WriteString("---\n")
			}
		}
		markdown.WriteString(post.OriginalContent)

		fileName := post.CreateTime.Format("2006-01-02") + "-" + post.Slug + ".md"
		file, err := os.OpenFile(filepath.Join(backupFilePath, fileName), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			return "", xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("create file err")
		}
		_, err = file.WriteString(markdown.String())
		if err != nil {
			return "", xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("write file err")
		}
		toBackupPaths = append(toBackupPaths, filepath.Join(backupFilePath, fileName))
	}

	backupFile := filepath.Join(backupFilePath, backupFilename)

	err = util.ZipFile(backupFile, toBackupPaths...)
	if err != nil {
		return "", err
	}
	return backupFile, nil
}

func (e *exportImport) getFrontMatterYaml(ctx context.Context, post *entity.Post) (string, error) {
	tags, err := e.PostTagService.ListTagByPostID(ctx, post.ID)
	if err != nil {
		return "", err
	}
	categories, err := e.PostCategoryService.ListCategoryByPostID(ctx, post.ID)
	if err != nil {
		return "", err
	}
	tagsStr := make([]string, 0)
	for _, tag := range tags {
		tagsStr = append(tagsStr, tag.Name)
	}
	categoriesStr := make([]string, 0)
	for _, category := range categories {
		categoriesStr = append(categoriesStr, category.Name)
	}
	frontMatter := make(map[string]any)
	frontMatter["title"] = post.Title
	frontMatter["draft"] = post.Status == consts.PostStatusDraft
	frontMatter["date"] = post.CreateTime.Format("2006-01-02 15:04:05 MST")
	frontMatter["comments"] = !post.DisallowComment
	frontMatter["slug"] = post.Slug
	if post.EditTime != nil {
		frontMatter["lastmod"] = post.EditTime.Format("2006-01-02 15:04:05 MST")
	}

	if post.UpdateTime != nil && post.UpdateTime != (&time.Time{}) {
		frontMatter["updated"] = post.UpdateTime.Format("2006-01-02 15:04:05 MST")
	}
	if post.Summary != "" {
		frontMatter["summary"] = post.Summary
	}
	if len(tagsStr) > 0 {
		frontMatter["tags"] = tagsStr
	}
	if len(categoriesStr) > 0 {
		frontMatter["categories"] = categoriesStr
	}
	out, err := yaml.Marshal(frontMatter)
	if err != nil {
		return "", xerr.WithStatus(err, xerr.StatusInternalServerError)
	}
	return string(out), nil
}

func convertJekyllMetaData(metadata map[string]any, postName string, postDate time.Time) map[string]any {
	for key, value := range metadata {
		lowerKey := strings.ToLower(key)

		switch lowerKey {
		case "layout":
			delete(metadata, key)
		case "permalink":
			if str, ok := value.(string); ok {
				metadata["url"] = str
			}
		case "category":
			if str, ok := value.(string); ok {
				metadata["categories"] = []string{str}
			}
			delete(metadata, key)
		case "excerpt_separator":
			if key != lowerKey {
				delete(metadata, key)
				metadata[lowerKey] = value
			}
		case "date":
			date, err := cast.StringToDate(value.(string))
			if err != nil {
				log.Errorf("convertJekyllMetaData date parse err date=%v err=%v", value.(string), err)
			} else {
				postDate = date
			}
		case "title":
			postName = value.(string)
		case "published":
			published, err := strconv.ParseBool(value.(string))
			if err != nil {
				log.Errorf("convertJekyllMetaData parse published err published=%v err=%v", value.(string), err)
			} else {
				delete(metadata, key)
				metadata["draft"] = strconv.FormatBool(published)
			}
		}
	}

	metadata["date"] = postDate.Format(time.RFC3339)
	metadata["title"] = postName
	return metadata
}

func convertJekyllContent(metadata map[string]any, content string) string {
	lines := strings.Split(content, "\n")
	resultLines := make([]string, 0, len(lines))
	for _, line := range lines {
		resultLines = append(resultLines, strings.Trim(line, "\r\n"))
	}

	content = strings.Join(resultLines, "\n")

	excerptSep := "<!--more-->"
	if value, ok := metadata["excerpt_separator"]; ok {
		if str, strOk := value.(string); strOk {
			content = strings.ReplaceAll(content, strings.TrimSpace(str), excerptSep)
		}
	}

	replaceList := []struct {
		re      *regexp.Regexp
		replace string
	}{
		{regexp.MustCompile("(?i)<!-- more -->"), "<!--more-->"},
		{regexp.MustCompile(`\{%\s*raw\s*%\}\s*(.*?)\s*\{%\s*endraw\s*%\}`), "$1"},
		{regexp.MustCompile(`{%\s*endhighlight\s*%}`), "{{< / highlight >}}"},
	}

	for _, replace := range replaceList {
		content = replace.re.ReplaceAllString(content, replace.replace)
	}

	replaceListFunc := []struct {
		re      *regexp.Regexp
		replace func(string) string
	}{
		// Octopress image tag: http://octopress.org/docs/plugins/image-tag/
		{regexp.MustCompile(`{%\s+img\s*(.*?)\s*%}`), replaceImageTag},
		{regexp.MustCompile(`{%\s*highlight\s*(.*?)\s*%}`), replaceHighlightTag},
	}

	for _, replace := range replaceListFunc {
		content = replace.re.ReplaceAllStringFunc(content, replace.replace)
	}

	return content
}

func replaceHighlightTag(match string) string {
	r := regexp.MustCompile(`{%\s*highlight\s*(.*?)\s*%}`)
	parts := r.FindStringSubmatch(match)
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}
	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(parts[1], f)

	result := bytes.NewBufferString("{{< highlight ")
	result.WriteString(items[0]) // language
	options := items[1:]
	for i, opt := range options {
		opt = strings.ReplaceAll(opt, "\"", "")
		if opt == "linenos" {
			opt = "linenos=table"
		}
		if i == 0 {
			opt = " \"" + opt
		}
		if i < len(options)-1 {
			opt += ","
		} else if i == len(options)-1 {
			opt += "\""
		}
		result.WriteString(opt)
	}

	result.WriteString(" >}}")
	return result.String()
}

func replaceImageTag(match string) string {
	r := regexp.MustCompile(`{%\s+img\s*(\p{L}*)\s+([\S]*/[\S]+)\s+(\d*)\s*(\d*)\s*(.*?)\s*%}`)
	result := bytes.NewBufferString("{{< figure ")
	parts := r.FindStringSubmatch(match)
	// Index 0 is the entire string, ignore
	replaceOptionalPart(result, "class", parts[1])
	replaceOptionalPart(result, "src", parts[2])
	replaceOptionalPart(result, "width", parts[3])
	replaceOptionalPart(result, "height", parts[4])
	// title + alt
	part := parts[5]
	if len(part) > 0 {
		splits := strings.Split(part, "'")
		lenSplits := len(splits)
		switch lenSplits {
		case 1:
			replaceOptionalPart(result, "title", splits[0])
		case 3:
			replaceOptionalPart(result, "title", splits[1])
		case 5:
			replaceOptionalPart(result, "title", splits[1])
			replaceOptionalPart(result, "alt", splits[3])
		}
	}
	result.WriteString(">}}")
	return result.String()
}

func replaceOptionalPart(buffer *bytes.Buffer, partName string, part string) {
	if len(part) > 0 {
		buffer.WriteString(partName + "=\"" + part + "\" ")
	}
}

func parseJekyllFilename(filename string) (time.Time, string, error) {
	re := regexp.MustCompile(`(\d+-\d+-\d+)-(.+)\..*`)
	r := re.FindAllStringSubmatch(filename, -1)
	if len(r) == 0 {
		return time.Now(), "", xerr.NoType.New("filename not match")
	}

	postDate, err := time.Parse("2006-1-2", r[0][1])
	if err != nil {
		return time.Now(), "", err
	}

	postName := r[0][2]

	return postDate, postName, nil
}
