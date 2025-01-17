package ks

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/yliu7949/KouShare-dl/live"
	"github.com/yliu7949/KouShare-dl/slide"
	"github.com/yliu7949/KouShare-dl/user"
	"github.com/yliu7949/KouShare-dl/video"
)

var path string

// InfoCmd 获取视频或直播的基本信息
func InfoCmd() *cobra.Command {
	var v video.Video
	var l live.Live
	var cmdInfo = &cobra.Command{
		Use:   "info [vid]",
		Short: "获取视频或直播的基本信息",
		Long:  `获取视频的基本信息，如讲者、拍摄日期、视频大小、视频摘要等内容；获取直播的基本信息，如开播时间、主办方、有无回放等内容.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args[0]) == 6 {
				l.RoomID = args[0]
				l.ShowLiveInfo()
			} else {
				v.Vid = args[0]
				v.ShowVideoInfo()
			}
		},
	}

	return cmdInfo
}

// SaveCmd 保存指定vid的视频
func SaveCmd() *cobra.Command {
	var v video.Video
	var quality string
	var isSeries bool
	var cmdSave = &cobra.Command{
		Use:   "save [vid]",
		Short: "保存指定vid的视频",
		Long:  `保存指定vid的视频到本地计算机，未登陆时仅可下载标清视频，登录后可以下载更高清晰度的视频.仅能下载已购买的付费视频.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			v.Vid = args[0]
			if path[len(path)-1:] != `\` && path[len(path)-1:] != "/" {
				path = path + "/"
			}
			v.SaveDir = path
			if isSeries {
				v.DownloadSeriesVideos(quality)
			} else {
				v.DownloadSingleVideo(quality)
			}
		},
		Aliases: []string{"video"},
	}
	cmdSave.Flags().StringVarP(&path, "path", "p", `.`, "指定保存视频的路径")
	cmdSave.Flags().BoolVarP(&isSeries, "series", "s", false, "指定是否下载专题视频")
	cmdSave.Flags().StringVarP(&quality, "quality", "q", `high`, "指定下载视频的清晰度（high、standard或low）")

	return cmdSave
}

// RecordCmd 录制指定直播间ID的直播
func RecordCmd() *cobra.Command {
	var l live.Live
	var liveTime string //开播时间，格式应为"2006-01-02 15:04:05"
	var autoMerge bool
	var replay bool

	var cmdRecord = &cobra.Command{
		Use:   "record [roomID]",
		Short: "录制指定直播间ID的直播",
		Long:  `录制指定直播间ID的直播.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			l.RoomID = args[0]
			if path[len(path)-1:] != `\` && path[len(path)-1:] != "/" {
				path = path + "/"
			}
			l.SaveDir = path
			if !replay {
				l.WaitAndRecordTheLive(liveTime, autoMerge)
			} else {
				l.DownloadReplayVideo()
			}
		},
		Aliases: []string{"live"},
	}
	cmdRecord.Flags().StringVarP(&path, "path", "p", `.`, "指定保存视频的路径")
	cmdRecord.Flags().StringVarP(&liveTime, "at", "@", "", `开播时间，格式为"2006-01-02 15:04:05"`)
	cmdRecord.Flags().BoolVarP(&autoMerge, "autoMerge", "a", false, "指定是否自动合并下载的视频片段文件")
	cmdRecord.Flags().BoolVarP(&replay, "replay", "r", false, "指定是否下载直播间快速回放视频")

	return cmdRecord
}

// MergeCmd 合并下载的视频片段文件
func MergeCmd() *cobra.Command {
	var dstFileName string
	var cmdMerge = &cobra.Command{
		Use:   "merge [directory]",
		Short: "合并下载的视频片段文件",
		Long:  `合并下载的视频片段文件(.ts)为一个视频文件(.ts)，[directory]参数为存放视频片段文件的文件夹的路径，若为空则默认为当前路径.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				path = "./"
			} else {
				path = args[0]
				if path[len(path)-1:] != `\` && path[len(path)-1:] != "/" {
					path = path + "/"
				}
			}
			live.MergeTsFiles(path, dstFileName)
		},
	}
	cmdMerge.Flags().StringVarP(&dstFileName, "name", "n", `recorded Video File.ts`, "指定合并后视频文件的名字(xxx.ts)")

	return cmdMerge
}

// SlideCmd 下载指定vid的视频对应的课件
func SlideCmd() *cobra.Command {
	var s slide.Slide
	var isSeries bool
	var qpdfBinPath string
	var cmdSlide = &cobra.Command{
		Use:   "slide [vid]",
		Short: "下载指定vid的视频对应的课件",
		Long:  `下载指定vid的视频对应的课件.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			s.Vid = args[0]
			if path[len(path)-1:] != `\` && path[len(path)-1:] != "/" {
				path = path + "/"
			}
			s.SaveDir = path
			if qpdfBinPath != "" {
				if qpdfBinPath[len(qpdfBinPath)-1:] != `\` && qpdfBinPath[len(qpdfBinPath)-1:] != "/" {
					qpdfBinPath = qpdfBinPath + "/"
				}
			}
			s.QpdfPath = qpdfBinPath
			if isSeries {
				s.DownloadSeriesSlides()
			} else {
				s.DownloadSingleSlide()
			}
		},
	}
	cmdSlide.Flags().StringVarP(&path, "path", "p", `.`, "指定保存课件的路径")
	cmdSlide.Flags().BoolVarP(&isSeries, "series", "s", false, "指定是否下载整个专题的所有课件")
	cmdSlide.Flags().StringVar(&qpdfBinPath, "qpdf-bin", "", "指定qpdf的bin文件夹所在的路径")

	return cmdSlide
}

// LoginCmd 通过短信验证码获取“蔻享学术”登陆凭证
func LoginCmd() *cobra.Command {
	var u user.User
	var cmdLogin = &cobra.Command{
		Use:   "login [phone number]",
		Short: "通过短信验证码获取“蔻享学术”登陆凭证",
		Long:  `[phone number]参数为手机号码（格式15012345678），输入短信验证码以登陆“蔻享学术”平台并将登陆凭证保存至本地.登录后一周内免再次登录.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			re := regexp.MustCompile(`1[3-9]\d{9}`)
			if !re.MatchString(args[0]) {
				fmt.Println("手机号码格式不正确")
				return
			}
			u.PhoneNumber = args[0]
			if err := u.Login(); err != nil {
				fmt.Println("登录失败：", err)
				return
			}
		},
	}

	return cmdLogin
}

// LogoutCmd 退出登录并删除保存在本地的登陆凭证文件
func LogoutCmd() *cobra.Command {
	var u user.User
	var cmdLogout = &cobra.Command{
		Use:   "logout",
		Short: "退出登陆",
		Long:  `退出登录并删除保存在本地的登陆凭证文件.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			u.Logout()
		},
	}

	return cmdLogout
}
