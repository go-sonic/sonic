package wp

type PostDTO struct {
	Date              string                 `json:"date"`
	DateGmt           string                 `json:"date_gmt"`
	GUID              map[string]interface{} `json:"guid"`
	ID                int32                  `json:"id"`
	Link              string                 `json:"link"`
	Modified          string                 `json:"modified"`
	ModifiedGmt       string                 `json:"modified_gmt"`
	Slug              string                 `json:"slug"`
	Status            string                 `json:"status"`
	Type              string                 `json:"type"`
	Password          string                 `json:"password"`
	PermalinkTemplate string                 `json:"permalink_template"`
	GeneratedSlug     string                 `json:"generated_slug"`
	Title             string                 `json:"title"`
	Content           map[string]interface{} `json:"content"`
	Author            int32                  `json:"author"`
	Excerpt           map[string]interface{} `json:"excerpt"`
	FeaturedMedia     int32                  `json:"featured_media"`
	CommentStatus     string                 `json:"comment_status"`
	PingStatus        string                 `json:"ping_status"`
	Format            string                 `json:"format"`
	Meta              map[string]interface{} `json:"meta"`
	Sticky            bool                   `json:"sticky"`
	Template          string                 `json:"template"`
	Categories        []int32                `json:"categories"`
	Tags              []int32                `json:"tags"`
}
