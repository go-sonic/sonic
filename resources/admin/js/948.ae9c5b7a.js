"use strict";(self["webpackChunkhalo_admin"]=self["webpackChunkhalo_admin"]||[]).push([[948],{53948:function(e,t,n){n.r(t),n.d(t,{default:function(){return p}});var a=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("page-view",{attrs:{"sub-title":e.theme.current.version||"-",title:e.theme.current.name||"-",affix:""}},[n("template",{slot:"extra"},[n("ReactiveButton",{attrs:{errored:e.form.saveErrored,loading:e.form.saving,erroredText:"保存失败",loadedText:"保存成功",icon:"save",text:"保存设置",type:"primary"},on:{callback:e.handleSaveSettingsCallback,click:e.handleSaveSettings}}),n("a-dropdown",[n("a-menu",{attrs:{slot:"overlay"},slot:"overlay"},[n("a-menu-item",{key:"1",on:{click:e.handleRemoteUpdate}},[n("a-icon",{attrs:{type:"cloud"}}),e._v(" 在线更新 ")],1),n("a-menu-item",{key:"2",on:{click:function(t){e.localUpgradeModel.visible=!0}}},[n("a-icon",{attrs:{type:"file"}}),e._v(" 本地更新 ")],1)],1),n("a-button",{attrs:{icon:"upload"}},[e._v(" 更新 "),n("a-icon",{attrs:{type:"down"}})],1)],1),n("a-dropdown",{scopedSlots:e._u([{key:"overlay",fn:function(){return[n("a-menu",[n("a-menu-item",{attrs:{disabled:e.theme.current.activated},on:{click:e.handleActiveTheme}},[n("a-icon",{attrs:{type:"lock"}}),e._v(" 启用 ")],1),n("a-menu-item",{attrs:{disabled:!e.theme.current.activated},on:{click:e.handleRouteToThemeVisualSetting}},[n("a-icon",{attrs:{type:"eye"}}),e._v(" 预览模式 ")],1)],1)]},proxy:!0}])},[n("a-button",{attrs:{icon:"more"}},[e._v(" 更多 "),n("a-icon",{attrs:{type:"down"}})],1)],1),n("a-button",{attrs:{disabled:e.theme.current.activated,icon:"delete",type:"danger"},on:{click:function(t){e.themeDeleteModal.visible=!0}}},[e._v(" 删除 ")])],1),n("a-spin",{attrs:{spinning:e.theme.loading}},[n("ThemeSettingForm",{ref:"themeSettingForm",attrs:{theme:e.theme.current}})],1),n("ThemeDeleteConfirmModal",{attrs:{theme:e.theme.current,visible:e.themeDeleteModal.visible},on:{"update:visible":function(t){return e.$set(e.themeDeleteModal,"visible",t)},success:e.onThemeDeleteSucceed}}),n("ThemeLocalUpgradeModal",{attrs:{theme:e.theme.current,visible:e.localUpgradeModel.visible},on:{"update:visible":function(t){return e.$set(e.localUpgradeModel,"visible",t)},success:e.handleGetTheme}})],2)},r=[],i=n(54288),s=(n(30535),n(85018),n(21082),n(70315),n(53298)),l=n(37130),o=n(10291),c=n(99868),u=n(36591),m={name:"ThemeSetting",components:{PageView:s.B4,ThemeDeleteConfirmModal:l.Z,ThemeLocalUpgradeModal:o.Z,ThemeSettingForm:c.Z},data:function(){return{theme:{current:{},loading:!1},themeDeleteModal:{visible:!1},localUpgradeModel:{visible:!1},form:{saving:!1,saveErrored:!1}}},beforeRouteEnter:function(e,t,n){var a=e.query.themeId;n(function(){var e=(0,i.Z)(regeneratorRuntime.mark((function e(t){return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.next=2,t.handleGetTheme(a);case 2:case"end":return e.stop()}}),e)})));return function(t){return e.apply(this,arguments)}}())},methods:{handleGetTheme:function(e){var t=this;return(0,i.Z)(regeneratorRuntime.mark((function n(){var a,r,i,s;return regeneratorRuntime.wrap((function(n){while(1)switch(n.prev=n.next){case 0:if(n.prev=0,t.theme.loading=!0,!e){n.next=10;break}return n.next=5,u.Z.theme.get(e);case 5:a=n.sent,r=a.data,t.theme.current=r,n.next=15;break;case 10:return n.next=12,u.Z.theme.getActivatedTheme();case 12:i=n.sent,s=i.data,t.theme.current=s;case 15:return n.prev=15,t.theme.loading=!1,n.finish(15);case 18:case"end":return n.stop()}}),n,null,[[0,,15,18]])})))()},onThemeDeleteSucceed:function(){this.$router.replace({name:"ThemeList"})},handleRemoteUpdate:function(){var e=this;e.$confirm({title:"提示",maskClosable:!0,content:"确定更新【"+e.theme.current.name+"】主题吗？",onOk:function(){return(0,i.Z)(regeneratorRuntime.mark((function t(){var n;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return n=e.$message.loading("更新中...",0),t.prev=1,t.next=4,u.Z.theme.updateThemeByFetching(e.theme.current.id);case 4:e.$message.success("更新成功！"),t.next=10;break;case 7:t.prev=7,t.t0=t["catch"](1),e.$log.error("Failed to update theme: ",t.t0);case 10:return t.prev=10,n(),t.next=14,e.handleGetTheme(e.theme.current.id);case 14:return t.finish(10);case 15:case"end":return t.stop()}}),t,null,[[1,7,10,15]])})))()}})},handleRouteToThemeVisualSetting:function(){this.$router.push({name:"ThemeVisualSetting",query:{themeId:this.theme.current.id}})},handleActiveTheme:function(){var e=this;e.$confirm({title:"提示",maskClosable:!0,content:"确定启用【"+e.theme.current.name+"】主题吗？",onOk:function(){return(0,i.Z)(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.prev=0,t.next=3,u.Z.theme.active(e.theme.current.id);case 3:e.$message.success("启用成功！"),t.next=9;break;case 6:t.prev=6,t.t0=t["catch"](0),e.$log.error("Failed active theme",t.t0);case 9:return t.prev=9,t.next=12,e.handleGetTheme(e.theme.current.id);case 12:return t.finish(9);case 13:case"end":return t.stop()}}),t,null,[[0,6,9,13]])})))()}})},handleSaveSettings:function(){var e=this;return(0,i.Z)(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.prev=0,e.form.saving=!0,t.next=4,e.$refs.themeSettingForm.handleSaveSettings(!1);case 4:t.next=9;break;case 6:t.prev=6,t.t0=t["catch"](0),e.form.saveErrored=!0;case 9:return t.prev=9,setTimeout((function(){e.form.saving=!1}),400),t.finish(9);case 12:case"end":return t.stop()}}),t,null,[[0,6,9,12]])})))()},handleSaveSettingsCallback:function(){this.form.saveErrored&&(this.form.saveErrored=!1)}}},d=m,h=n(70739),f=(0,h.Z)(d,a,r,!1,null,null,null),p=f.exports},37130:function(e,t,n){n.d(t,{Z:function(){return m}});var a=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("a-modal",{attrs:{afterClose:e.onAfterClose,closable:!1,width:416,destroyOnClose:"",title:"提示"},model:{value:e.modalVisible,callback:function(t){e.modalVisible=t},expression:"modalVisible"}},[n("template",{slot:"footer"},[n("a-button",{on:{click:function(t){e.modalVisible=!1}}},[e._v(" 取消 ")]),n("ReactiveButton",{attrs:{errored:e.deleteErrored,loading:e.deleting,erroredText:"删除失败",loadedText:"删除成功",text:"确定"},on:{callback:e.handleDeleteCallback,click:function(t){return e.handleDelete()}}})],1),n("p",[e._v("确定删除【"+e._s(e.theme.name)+"】主题？")]),n("a-checkbox",{model:{value:e.deleteSettings,callback:function(t){e.deleteSettings=t},expression:"deleteSettings"}},[e._v(" 同时删除主题配置 ")])],2)},r=[],i=n(54288),s=(n(70315),n(36591)),l={name:"ThemeDeleteConfirmModal",props:{visible:{type:Boolean,default:!1},theme:{type:Object,default:function(){return{}}}},data:function(){return{deleteErrored:!1,deleting:!1,deleteSettings:!1}},computed:{modalVisible:{get:function(){return this.visible},set:function(e){this.$emit("update:visible",e)}}},methods:{handleDelete:function(){var e=this;return(0,i.Z)(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.prev=0,e.deleting=!0,t.next=4,s.Z.theme["delete"](e.theme.id,e.deleteSettings);case 4:t.next=10;break;case 6:t.prev=6,t.t0=t["catch"](0),e.deleteErrored=!1,e.$log.error("Delete theme failed",t.t0);case 10:return t.prev=10,setTimeout((function(){e.deleting=!1}),400),t.finish(10);case 13:case"end":return t.stop()}}),t,null,[[0,6,10,13]])})))()},handleDeleteCallback:function(){this.deleteErrored?this.deleteErrored=!1:(this.modalVisible=!1,this.$emit("success"))},onAfterClose:function(){this.deleteErrored=!1,this.deleting=!1,this.deleteSettings=!1,this.$emit("onAfterClose")}}},o=l,c=n(70739),u=(0,c.Z)(o,a,r,!1,null,null,null),m=u.exports},10291:function(e,t,n){n.d(t,{Z:function(){return u}});var a=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("a-modal",{attrs:{afterClose:e.onModalClose,footer:null,destroyOnClose:"",title:"更新主题"},model:{value:e.modalVisible,callback:function(t){e.modalVisible=t},expression:"modalVisible"}},[n("FilePondUpload",{ref:"updateByFile",attrs:{accepts:["application/x-zip","application/x-zip-compressed","application/zip"],field:e.theme.id,multiple:!1,uploadHandler:e.uploadHandler,label:"点击选择主题更新包或将主题更新包拖拽到此处<br>仅支持 ZIP 格式的文件",name:"file"},on:{success:e.onThemeUploadSuccess}})],1)},r=[],i=n(36591),s={name:"ThemeLocalUpgradeModal",props:{visible:{type:Boolean,default:!1},theme:{type:Object,default:function(){return{}}}},data:function(){return{uploadHandler:function(e,t,n){return i.Z.theme.updateByUpload(e,t,n)}}},computed:{modalVisible:{get:function(){return this.visible},set:function(e){this.$emit("update:visible",e)}}},methods:{onModalClose:function(){this.$refs.updateByFile.handleClearFileList(),this.$emit("onAfterClose")},onThemeUploadSuccess:function(){this.modalVisible=!1,this.$emit("success")}}},l=s,o=n(70739),c=(0,o.Z)(l,a,r,!1,null,null,null),u=c.exports},99868:function(e,t,n){n.d(t,{Z:function(){return d}});var a=function(){var e=this,t=e.$createElement,n=e._self._c||t;return e.theme.id?n("div",{staticClass:"card-container h-full"},[n("a-tabs",{staticClass:"h-full",attrs:{defaultActiveKey:"0",type:"card"}},[n("a-tab-pane",{key:0,attrs:{tab:"关于"}},[e.theme.logo?n("div",[n("a-avatar",{attrs:{alt:e.theme.name,size:72,src:e.theme.logo,shape:"square"}}),n("a-divider")],1):e._e(),n("a-descriptions",{attrs:{column:1,layout:"horizontal"}},[n("a-descriptions-item",{attrs:{label:"作者"}},[n("a",{staticClass:"text-inherit",attrs:{href:e.theme.author.website||"#",target:"_blank"}},[e._v(" "+e._s(e.theme.author.name)+" ")])]),n("a-descriptions-item",{attrs:{label:"介绍"}},[e._v(" "+e._s(e.theme.description||"-")+" ")]),n("a-descriptions-item",{attrs:{label:"官网"}},[n("a",{staticClass:"text-inherit",attrs:{href:e.theme.website||"#",target:"_blank"}},[e._v(" "+e._s(e.theme.website||"-")+" ")])]),n("a-descriptions-item",{attrs:{label:"Git 仓库"}},[n("a",{staticClass:"text-inherit",attrs:{href:e.theme.repo||"#",target:"_blank"}},[e._v(" "+e._s(e.theme.repo||"-")+" ")])]),n("a-descriptions-item",{attrs:{label:"主题标识"}},[e._v(" "+e._s(e.theme.id)+" ")]),n("a-descriptions-item",{attrs:{label:"当前版本"}},[e._v(" "+e._s(e.theme.version)+" ")]),n("a-descriptions-item",{attrs:{label:"存储位置"}},[e._v(" "+e._s(e.theme.themePath)+" ")]),e._t("descriptions-item")],2)],1),e._l(e.form.configurations,(function(t,a){return n("a-tab-pane",{key:a+1,attrs:{tab:t.label}},[n("a-form",{attrs:{wrapperCol:e.wrapperCol,layout:"vertical"}},[e._l(t.items,(function(t,a){return n("a-form-item",{key:a,attrs:{label:t.label+"："}},[t.description&&""!==t.description?n("p",{attrs:{slot:"help"},domProps:{innerHTML:e._s(t.description)},slot:"help"}):e._e(),"TEXT"===t.type?n("a-input",{attrs:{defaultValue:t.defaultValue,placeholder:t.placeholder},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}}):"TEXTAREA"===t.type?n("a-input",{attrs:{autoSize:{minRows:5},placeholder:t.placeholder,type:"textarea"},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}}):"RADIO"===t.type?n("a-radio-group",{attrs:{defaultValue:t.defaultValue},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}},e._l(t.options,(function(t,a){return n("a-radio",{key:a,attrs:{value:t.value}},[e._v(" "+e._s(t.label)+" ")])})),1):"SELECT"===t.type?n("a-select",{attrs:{defaultValue:t.defaultValue},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}},e._l(t.options,(function(t){return n("a-select-option",{key:t.value,attrs:{value:t.value}},[e._v(" "+e._s(t.label)+" ")])})),1):"COLOR"===t.type?n("verte",{staticStyle:{display:"inline-block",height:"24px"},attrs:{defaultValue:t.defaultValue,model:"hex",picker:"square"},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}}):"ATTACHMENT"===t.type?n("AttachmentInput",{attrs:{defaultValue:t.defaultValue,placeholder:t.placeholder},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}}):"NUMBER"===t.type?n("a-input-number",{staticStyle:{width:"100%"},attrs:{defaultValue:t.defaultValue},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}}):"SWITCH"===t.type?n("a-switch",{attrs:{defaultChecked:t.defaultValue},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}}):n("a-input",{attrs:{defaultValue:t.defaultValue,placeholder:t.placeholder},model:{value:e.form.settings[t.name],callback:function(n){e.$set(e.form.settings,t.name,n)},expression:"form.settings[item.name]"}})],1)})),n("a-form-item",[n("ReactiveButton",{attrs:{errored:e.form.saveErrored,loading:e.form.saving,erroredText:"保存失败",loadedText:"保存成功",text:"保存",type:"primary"},on:{callback:e.handleSaveSettingsCallback,click:e.handleSaveSettings}})],1)],2)],1)}))],2)],1):e._e()},r=[],i=n(54288),s=(n(87591),n(70315),n(43154)),l=n(36591),o={name:"ThemeSettingForm",components:{Verte:s.Z},props:{theme:{type:Object,default:function(){}},wrapperCol:{type:Object,default:function(){return{xl:{span:8},lg:{span:8},sm:{span:12},xs:{span:24}}}}},data:function(){return{form:{settings:[],configurations:[],loading:!1,saving:!1,saveErrored:!1}}},watch:{theme:function(e){e&&this.handleGetConfigurations()}},methods:{handleGetConfigurations:function(){var e=this;return(0,i.Z)(regeneratorRuntime.mark((function t(){var n,a;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.prev=0,t.next=3,l.Z.theme.listConfigurations(e.theme.id);case 3:return n=t.sent,a=n.data,e.form.configurations=a,t.next=8,e.handleGetSettings();case 8:t.next=13;break;case 10:t.prev=10,t.t0=t["catch"](0),e.$log.error(t.t0);case 13:case"end":return t.stop()}}),t,null,[[0,10]])})))()},handleGetSettings:function(){var e=this;return(0,i.Z)(regeneratorRuntime.mark((function t(){var n,a;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.prev=0,t.next=3,l.Z.theme.listSettings(e.theme.id);case 3:n=t.sent,a=n.data,e.form.settings=a,t.next=11;break;case 8:t.prev=8,t.t0=t["catch"](0),e.$log.error(t.t0);case 11:case"end":return t.stop()}}),t,null,[[0,8]])})))()},handleSaveSettings:function(){var e=arguments,t=this;return(0,i.Z)(regeneratorRuntime.mark((function n(){var a;return regeneratorRuntime.wrap((function(n){while(1)switch(n.prev=n.next){case 0:return a=!(e.length>0&&void 0!==e[0])||e[0],n.prev=1,a&&(t.form.saving=!0),n.next=5,l.Z.theme.saveSettings(t.theme.id,t.form.settings);case 5:n.next=12;break;case 7:throw n.prev=7,n.t0=n["catch"](1),t.$log.error(n.t0),t.form.saveErrored=!0,new Error(n.t0);case 12:return n.prev=12,setTimeout((function(){t.form.saving=!1}),400),n.finish(12);case 15:case"end":return n.stop()}}),n,null,[[1,7,12,15]])})))()},handleSaveSettingsCallback:function(){this.form.saveErrored?this.form.saveErrored=!1:(this.handleGetSettings(),this.$emit("saved"))}}},c=o,u=n(70739),m=(0,u.Z)(c,a,r,!1,null,null,null),d=m.exports}}]);