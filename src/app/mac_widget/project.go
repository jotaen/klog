package mac_widget

//
//import (
//	"github.com/caseymrm/menuet"
//	"klog/app"
//	"klog/datetime"
//	"klog/project"
//	"klog/record"
//	"time"
//)
//
//func renderProject(project project.Project) []menuet.MenuItem {
//
//	items = append(items, menuet.MenuItem{
//		Text: "History",
//		Children: func() []menuet.MenuItem {
//			items := []menuet.MenuItem{
//				{Text: "Open folder", Clicked: func() {
//					app.OpenInFileBrowser(project)
//				}},
//			}
//			days, _ := project.List()
//			if len(days) > 0 {
//				items = append(items, menuet.MenuItem{Type: menuet.Separator})
//				for i, d := range days {
//					if i == 5 {
//						break
//					}
//					wd, _ := project.Get(d)
//					if wd == nil {
//						break
//					}
//					todayLabel := ""
//					if currentDay != nil && wd.Date() == currentDay.Date() {
//						todayLabel = " (today)"
//					}
//					items = append(items, menuet.MenuItem{
//						Text: wd.Date().ToString() + todayLabel, Clicked: func() {
//							app.OpenInEditor(project, wd)
//						},
//					})
//				}
//			}
//			return items
//		},
//	})
//
//	return items
//}
