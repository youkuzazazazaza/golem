package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sbgayhub/golem/sdk/cdn"
	"github.com/sbgayhub/golem/sdk/contact"
	"github.com/sbgayhub/golem/sdk/message"
	"github.com/sbgayhub/golem/sdk/plugin"
)

// Config 插件配置
type Config struct {
	VideoNative bool   `toml:"video_native" comment:"视频使用 CDN 原生上传（失败则回退到链接卡片）"`
	MaxList     int    `toml:"max_list" comment:"列表类结果最大条数"`
	BdbkURL     string `toml:"bdbk_url" comment:"百度百科 API 地址（需自行配置）"`
}

// DemosPlugin 娱乐功能插件
type DemosPlugin struct {
	plugin.ConfigAbility[Config]
	message message.Ability
	contact contact.Ability
	cdn     cdn.Ability
	client  *http.Client

	handlers map[string]handlerFunc
}

type handlerFunc func(receiver *contact.Contact, arg string) (bool, error)

const (
	xjjURL      = "https://api.yujn.cn/api/zzxjj.php?type=video"
	xjj2URL     = "https://api.kuleu.com/api/MP4_xiaojiejie?type=json"
	rdVideoURL  = "https://api.52vmy.cn/api/video/redian"
	ylVideoURL  = "https://api.52vmy.cn/api/video/yule"
	boyURL      = "https://api.52vmy.cn/api/video/boy"
	sdjURL      = "https://api.kuleu.com/api/action?text="
	catURL      = "https://api.thecatapi.com/v1/images/search?limit=1"
	dogURL      = "https://dog.ceo/api/breeds/image/random"
	twURL       = "https://api.nasa.gov/planetary/apod?api_key=TJTjotiNFKFh541VXfSwmsKdwMBVuRUikDmyPCgN&count=1"
	paintingURL = "https://api.52vmy.cn/api/query/painting"
	wxtsURL     = "https://api.kuleu.com/api/getGreetingMessage?type=json"
	yijuURL     = "https://api.apiopen.top/api/tools/famous-sentence"
	shiciURL    = "https://v2.alapi.cn/api/shici?type=all&token=iildXgwOPO6d7BOa"
	hahaURL     = "https://v2.alapi.cn/api/joke/random?token=iildXgwOPO6d7BOa"
	jzwURL      = "https://api.52vmy.cn/api/wl/s/jzw"
	raoURL      = "https://api.52vmy.cn/api/wl/yan/rao"
	hunyanURL   = "https://api.52vmy.cn/api/img/tw/card"
	dogDocURL   = "https://api.52vmy.cn/api/wl/s/dog?msg="
	yanyuURL    = "https://api.52vmy.cn/api/wl/yan/yanyu"
	chouqianURL = "https://api.52vmy.cn/api/wl/s/draw"
	bayURL      = "https://api.52vmy.cn/api/wl/yan/bay"
	eatURL      = "https://api.52vmy.cn/api/wl/s/eat"
	kingTcURL   = "https://api.yujn.cn/api/wzry.php?type=json"
	lolTcURL    = "https://api.yujn.cn/api/yxlm.php?"
	xhzBqURL    = "http://api.yujn.cn/api/cxk.php?"
	sjecyURL    = "https://api.cenguigui.cn/api/pic/"
	rjURL       = "https://api.yujn.cn/api/baoan.php?"
	acgURL      = "https://api.yujn.cn/api/gzl_ACG.php?type=image&form=pc"

	lyyKgURL      = "http://api.yujn.cn/api/lyy.php?type=video"
	duilianURL    = "https://api.yujn.cn/api/duilian.php?type=video"
	chuandaURL    = "http://api.yujn.cn/api/chuanda.php?type=video"
	shwdURL       = "http://api.yujn.cn/api/shwd.php?type=video"
	ksFcURL       = "http://api.yujn.cn/api/ks_fc.php?type=video"
	sjkkURL       = "http://api.yujn.cn/api/sjkk.php?"
	sjSingURL     = "https://www.hhlqilongzhu.cn/api/changya.php"
	kuwoSearchURL = "http://search.kuwo.cn/r.s"
	kuwoMusicURL  = "http://nmobi.kuwo.cn/mobi.s"
	kuwoInfoURL   = "http://m.kuwo.cn/newh5/singles/songinfoandlrc"
)

var noRedirectClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("过多重定向")
			}
			loc := req.Response.Header.Get("Location")
			if loc != "" {
				loc = strings.Trim(loc, "'\"")
				if parsed, err := url.Parse(loc); err == nil {
					req.URL = req.URL.ResolveReference(parsed)
				}
			}
			return nil
		},
	}
}

func (p *DemosPlugin) GetMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Name:        "demos",
		Author:      "Golem Team",
		Version:     "1.0.0",
		Description: "Demos 娱乐功能插件",
		Priority:    0,
		Next:        false,
		AlwaysRun:   false,
	}
}

func (p *DemosPlugin) OnLoad() error {
	p.ensureDefaults()
	slog.Info("[demos] 插件加载成功", "video_native", p.Config.VideoNative, "max_list", p.Config.MaxList)
	return nil
}

func (p *DemosPlugin) OnUnload() error {
	slog.Info("[demos] 插件已卸载")
	return nil
}

func (p *DemosPlugin) GetSubscriptions() []string {
	return []string{message.TypeText.Topic}
}

func (p *DemosPlugin) ensureDefaults() {
	if p.client == nil {
		p.client = newHTTPClient()
	}
	if p.Config.MaxList == 0 {
		p.Config.MaxList = 3
	}
}

func (p *DemosPlugin) sortedKeys() []string {
	keys := make([]string, 0, len(p.handlers))
	for k := range p.handlers {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) > len(keys[j])
		}
		return keys[i] > keys[j]
	})
	return keys
}

func (p *DemosPlugin) OnEvent(e *plugin.Event) (bool, error) {
	p.ensureDefaults()

	msg := e.Payload.(*plugin.Event_Message).Message
	if msg == nil {
		return false, nil
	}

	text := strings.TrimSpace(msg.GetContent())
	if text == "" {
		return false, nil
	}

	receiver := p.contact.Get(e.GetSender())
	if receiver == nil {
		slog.Warn("[demos] 未找到接收者", "sender", e.GetSender())
		return false, nil
	}

	for _, key := range p.sortedKeys() {
		var arg string
		if text == key {
			arg = ""
		} else if strings.HasPrefix(text, key+" ") {
			arg = strings.TrimSpace(text[len(key):])
		} else {
			continue
		}

		handled, err := p.handlers[key](receiver, arg)
		if err != nil {
			slog.Error("[demos] 处理命令失败", "key", key, "err", err)
			p.sendText(receiver, "哎呀，翻车了！这个功能暂时罢工了，请稍后再试试吧~")
			return true, nil
		}
		return handled, nil
	}

	return false, nil
}

// ==================== 工具方法 ====================

func (p *DemosPlugin) httpGet(urlStr string) (string, error) {
	resp, err := p.client.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

func (p *DemosPlugin) httpGetWithHeaders(urlStr string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		if strings.EqualFold(k, "Host") {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

func (p *DemosPlugin) downloadMedia(urlStr string) ([]byte, error) {
	resp, err := p.client.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (p *DemosPlugin) getRedirectURL(u string) (string, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	resp, err := noRedirectClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", fmt.Errorf("没有 Location 头")
	}
	loc = strings.Trim(loc, "'\"")
	if strings.HasPrefix(loc, "http") {
		return loc, nil
	}
	base, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	rel, err := url.Parse(loc)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(rel).String(), nil
}

func (p *DemosPlugin) sendText(receiver *contact.Contact, text string) {
	msg := &message.Message{
		Type:     message.TypeText,
		Receiver: receiver,
		Content:  text,
		Data:     &message.Message_Text{Text: &message.TextData{Content: text}},
	}
	if _, err := p.message.Send(msg); err != nil {
		slog.Error("[demos] 发送文本失败", "err", err)
	}
}

func (p *DemosPlugin) sendImage(receiver *contact.Contact, imageURL string) error {
	data, err := p.downloadMedia(imageURL)
	if err != nil {
		return err
	}
	_, err = p.cdn.UploadImage(receiver.GetUsername(), bytes.NewReader(data))
	if err != nil {
		slog.Error("[demos] CDN 上传图片失败", "err", err)
		p.sendText(receiver, "图片发送失败，直接看链接吧："+imageURL)
		return nil
	}
	return nil
}

func (p *DemosPlugin) sendVideoCard(receiver *contact.Contact, title, desc, videoURL string) {
	xml := fmt.Sprintf(
		`<msg><appmsg appid="" sdkver="0"><title>%s</title><des>%s</des><action>view</action><type>5</type><showtype>0</showtype><url>%s</url><thumburl>%s</thumburl></appmsg></msg>`,
		escapeXML(title), escapeXML(desc), escapeXML(videoURL), escapeXML(defaultThumb()),
	)
	msg := &message.Message{
		Type:     message.TypeAppLink,
		Receiver: receiver,
		Content:  fmt.Sprintf("%s %s", title, desc),
		Data: &message.Message_App{App: &message.AppData{
			SubType: 5,
			Title:   title,
			Desc:    desc,
			Url:     videoURL,
			Xml:     xml,
		}},
	}
	if _, err := p.message.Send(msg); err != nil {
		slog.Error("[demos] 发送视频卡片失败", "err", err)
	}
}

func (p *DemosPlugin) sendVideoOrCard(receiver *contact.Contact, videoURL string) {
	if p.Config.VideoNative {
		err := p.sendNativeVideo(receiver, videoURL)
		if err == nil {
			return
		}
		slog.Warn("[demos] 原生视频发送失败，使用链接卡片", "err", err)
	}
	p.sendVideoCard(receiver, "视频链接", "点击播放", videoURL)
}

func (p *DemosPlugin) sendNativeVideo(receiver *contact.Contact, videoURL string) error {
	tmpVideo, err := os.CreateTemp("", "demos-video-*.mp4")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer os.Remove(tmpVideo.Name())
	defer tmpVideo.Close()

	resp, err := p.client.Get(videoURL)
	if err != nil {
		return fmt.Errorf("下载视频失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if _, err := io.Copy(tmpVideo, resp.Body); err != nil {
		return fmt.Errorf("保存视频失败: %w", err)
	}
	_ = tmpVideo.Close()

	duration, err := p.mediaDuration(tmpVideo.Name())
	if err != nil {
		slog.Warn("[demos] 获取视频时长失败，使用默认值", "err", err)
		duration = 10
	}

	thumbData, err := p.extractThumbnail(tmpVideo.Name())
	if err != nil {
		slog.Warn("[demos] 提取缩略图失败，使用空缩略图", "err", err)
		thumbData = []byte{}
	}

	videoFile, err := os.Open(tmpVideo.Name())
	if err != nil {
		return fmt.Errorf("打开视频文件失败: %w", err)
	}
	defer videoFile.Close()

	_, err = p.cdn.UploadVideo(receiver.GetUsername(), thumbData, videoFile, uint32(duration))
	if err != nil {
		return fmt.Errorf("CDN 上传视频失败: %w", err)
	}
	return nil
}

func (p *DemosPlugin) sendVoice(receiver *contact.Contact, audioURL string) {
	data, err := p.downloadMedia(audioURL)
	if err != nil {
		p.sendText(receiver, "语音获取失败，请稍后再试")
		return
	}

	tmpFile, err := os.CreateTemp("", "demos-audio-*.mp3")
	if err != nil {
		p.sendText(receiver, "语音处理失败")
		return
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write(data); err != nil {
		p.sendText(receiver, "语音处理失败")
		return
	}
	_ = tmpFile.Close()

	durationSec, err := p.mediaDuration(tmpFile.Name())
	if err != nil {
		durationSec = 0
	}

	msg := &message.Message{
		Type:     message.TypeVoice,
		Receiver: receiver,
		Content:  "[语音]",
		Data: &message.Message_Voice{Voice: &message.VoiceData{
			Media:    &message.Media{Data: data},
			Duration: uint32(durationSec * 1000),
		}},
	}
	if _, err := p.message.Send(msg); err != nil {
		slog.Error("[demos] 发送语音失败", "err", err)
		p.sendText(receiver, "语音发送失败，链接："+audioURL)
	}
}

func (p *DemosPlugin) mediaDuration(path string) (int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	d, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, err
	}
	return int(d), nil
}

func (p *DemosPlugin) extractThumbnail(videoPath string) ([]byte, error) {
	tmpThumb, err := os.CreateTemp("", "demos-thumb-*.jpg")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpThumb.Name())
	_ = tmpThumb.Close()

	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-f", "image2",
		"-y",
		tmpThumb.Name(),
	)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return os.ReadFile(tmpThumb.Name())
}

func defaultThumb() string {
	return "https://img0.baidu.com/it/u=3879589492,1588221464&fm=253&fmt=auto&app=120&f=JPEG?w=500&h=500"
}

func escapeXML(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	).Replace(s)
}

// ==================== 图片类处理器 ====================

func (p *DemosPlugin) handleCat(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(catURL)
	if err != nil {
		return true, err
	}
	var result []struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil || len(result) == 0 {
		return true, fmt.Errorf("解析猫咪图片失败")
	}
	return true, p.sendImage(receiver, result[0].URL)
}

func (p *DemosPlugin) handleDog(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(dogURL)
	if err != nil {
		return true, err
	}
	var result struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil || result.Message == "" {
		return true, fmt.Errorf("解析狗狗图片失败")
	}
	return true, p.sendImage(receiver, result.Message)
}

func (p *DemosPlugin) handleTw(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(twURL)
	if err != nil {
		return true, err
	}
	var list []struct {
		URL         string `json:"url"`
		Title       string `json:"title"`
		Date        string `json:"date"`
		Explanation string `json:"explanation"`
	}
	if err := json.Unmarshal([]byte(body), &list); err != nil || len(list) == 0 {
		return true, fmt.Errorf("解析天文图片失败")
	}
	data := list[0]
	p.sendText(receiver, fmt.Sprintf("看星空：%s\n时间：%s\n描述：%s", data.Title, data.Date, data.Explanation))
	return true, p.sendImage(receiver, data.URL)
}

func (p *DemosPlugin) handlePainting(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(paintingURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Img     string `json:"img"`
			Title   string `json:"title"`
			Dynasty string `json:"dynasty"`
			Source  string `json:"source"`
			Info    string `json:"info"`
			Content string `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析名画失败")
	}
	d := resp.Data
	p.sendText(receiver, fmt.Sprintf("《%s》\n--%s  %s\n%s\n%s", d.Title, d.Dynasty, d.Source, d.Info, d.Content))
	return true, p.sendImage(receiver, d.Img)
}

func (p *DemosPlugin) handleXhzbq(receiver *contact.Contact, arg string) (bool, error) {
	return true, p.sendImage(receiver, xhzBqURL)
}

func (p *DemosPlugin) handleSjecy(receiver *contact.Contact, arg string) (bool, error) {
	u, err := p.getRedirectURL(sjecyURL)
	if err != nil {
		return true, err
	}
	return true, p.sendImage(receiver, u)
}

func (p *DemosPlugin) handleAcg(receiver *contact.Contact, arg string) (bool, error) {
	u, err := p.getRedirectURL(acgURL)
	if err != nil {
		return true, err
	}
	return true, p.sendImage(receiver, u)
}

// ==================== 视频类处理器 ====================

func (p *DemosPlugin) handleXjj(receiver *contact.Contact, arg string) (bool, error) {
	u, err := p.getRedirectURL(xjjURL)
	if err != nil {
		return true, err
	}
	p.sendVideoOrCard(receiver, u)
	return true, nil
}

func (p *DemosPlugin) handleXjj2(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(xjj2URL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Mp4Video string `json:"mp4_video"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Mp4Video == "" {
		return true, fmt.Errorf("解析小姐姐视频失败")
	}
	p.sendVideoOrCard(receiver, resp.Mp4Video)
	return true, nil
}

func (p *DemosPlugin) handleRdVideo(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(rdVideoURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Data struct {
			Video string `json:"video"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Data.Video == "" {
		return true, fmt.Errorf("解析热点视频失败")
	}
	p.sendVideoOrCard(receiver, resp.Data.Video)
	return true, nil
}

func (p *DemosPlugin) handleYlVideo(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(ylVideoURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Data struct {
			Video string `json:"video"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Data.Video == "" {
		return true, fmt.Errorf("解析娱乐视频失败")
	}
	p.sendVideoOrCard(receiver, resp.Data.Video)
	return true, nil
}

func (p *DemosPlugin) handleBoyVideo(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(boyURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Data struct {
			Video string `json:"video"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Data.Video == "" {
		return true, fmt.Errorf("解析靓仔视频失败")
	}
	p.sendVideoOrCard(receiver, resp.Data.Video)
	return true, nil
}

func (p *DemosPlugin) handleLyyKg(receiver *contact.Contact, arg string) (bool, error) {
	u, err := p.getRedirectURL(lyyKgURL)
	if err != nil {
		return true, err
	}
	p.sendVideoOrCard(receiver, u)
	return true, nil
}

func (p *DemosPlugin) handleDuilian(receiver *contact.Contact, arg string) (bool, error) {
	u, err := p.getRedirectURL(duilianURL)
	if err != nil {
		return true, err
	}
	p.sendVideoOrCard(receiver, u)
	return true, nil
}

func (p *DemosPlugin) handleChuanda(receiver *contact.Contact, arg string) (bool, error) {
	u, err := p.getRedirectURL(chuandaURL)
	if err != nil {
		return true, err
	}
	p.sendVideoOrCard(receiver, u)
	return true, nil
}

func (p *DemosPlugin) handleShwd(receiver *contact.Contact, arg string) (bool, error) {
	u, err := p.getRedirectURL(shwdURL)
	if err != nil {
		return true, err
	}
	p.sendVideoOrCard(receiver, u)
	return true, nil
}

func (p *DemosPlugin) handleKsfc(receiver *contact.Contact, arg string) (bool, error) {
	u, err := p.getRedirectURL(ksFcURL)
	if err != nil {
		return true, err
	}
	p.sendVideoOrCard(receiver, u)
	return true, nil
}

// ==================== 文本类处理器 ====================

func (p *DemosPlugin) handleWxts(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(wxtsURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Greeting string `json:"greeting"`
			Tip      string `json:"tip"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析温馨提示失败")
	}
	p.sendText(receiver, fmt.Sprintf("%s\n%s", resp.Data.Greeting, resp.Data.Tip))
	return true, nil
}

func (p *DemosPlugin) handleYiju(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(yijuURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Name string `json:"name"`
			From string `json:"from"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析一句失败")
	}
	p.sendText(receiver, fmt.Sprintf("【%s】\n——【%s】", resp.Data.Name, resp.Data.From))
	return true, nil
}

func (p *DemosPlugin) handleYiyan(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet("https://v1.hitokoto.cn/")
	if err != nil {
		return true, err
	}
	var resp struct {
		Hitokoto string `json:"hitokoto"`
		From     string `json:"from"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Hitokoto == "" {
		return true, fmt.Errorf("解析一言失败")
	}
	if resp.From == "" {
		resp.From = "未知"
	}
	p.sendText(receiver, fmt.Sprintf("【%s】\n——【%s】", resp.Hitokoto, resp.From))
	return true, nil
}

func (p *DemosPlugin) handleShici(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(shiciURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Origin   string `json:"origin"`
			Author   string `json:"author"`
			Content  string `json:"content"`
			Category string `json:"category"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析诗词失败")
	}
	d := resp.Data
	p.sendText(receiver, fmt.Sprintf("【%s】\n——%s\n%s\n\n诗词类型：%s", d.Origin, d.Author, d.Content, d.Category))
	return true, nil
}

func (p *DemosPlugin) handleHaha(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(hahaURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析笑话失败")
	}
	p.sendText(receiver, fmt.Sprintf("\"%s\"\n%s", resp.Data.Title, resp.Data.Content))
	return true, nil
}

func (p *DemosPlugin) handleJzw(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(jzwURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Question string `json:"question"`
			Answer   string `json:"answer"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析脑筋急转弯失败")
	}
	p.sendText(receiver, fmt.Sprintf("问题：%s\n\n答案：%s", resp.Data.Question, resp.Data.Answer))
	return true, nil
}

func (p *DemosPlugin) handleRao(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(raoURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Title string `json:"title"`
			Msg   string `json:"msg"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析绕口令失败")
	}
	p.sendText(receiver, fmt.Sprintf("\"%s\"\n%s", resp.Data.Title, resp.Data.Msg))
	return true, nil
}

func (p *DemosPlugin) handleYanyu(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(yanyuURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Content string `json:"content"`
			Source  string `json:"source"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析谚语失败")
	}
	p.sendText(receiver, fmt.Sprintf("\"%s\"\n分类：%s", resp.Data.Content, resp.Data.Source))
	return true, nil
}

func (p *DemosPlugin) handleChouqian(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(chouqianURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析抽签失败")
	}
	p.sendText(receiver, fmt.Sprintf("\"%s\"\n%s", resp.Data.Text, resp.Data.ID))
	return true, nil
}

func (p *DemosPlugin) handleBay(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(bayURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Zh string `json:"zh"`
			En string `json:"en"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析答案之书失败")
	}
	p.sendText(receiver, fmt.Sprintf("%s\n%s", resp.Data.Zh, resp.Data.En))
	return true, nil
}

func (p *DemosPlugin) handleEat(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(eatURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析吃什么失败")
	}
	p.sendText(receiver, resp.Data)
	return true, nil
}

func (p *DemosPlugin) handleRj(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(rjURL)
	if err != nil {
		return true, err
	}
	body = strings.ReplaceAll(body, "\\n", "\n")
	body = strings.ReplaceAll(body, "\\t", "\t")
	p.sendText(receiver, body)
	return true, nil
}

func (p *DemosPlugin) handleKingTc(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(kingTcURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Name    string `json:"name"`
			Content string `json:"content"`
			Img     string `json:"img"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析 king 台词失败")
	}
	p.sendText(receiver, fmt.Sprintf("~%s~\n%s", resp.Data.Name, resp.Data.Content))
	return true, p.sendImage(receiver, resp.Data.Img)
}

func (p *DemosPlugin) handleLolTc(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(lolTcURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Name    string `json:"name"`
			Content string `json:"content"`
			Img     string `json:"img"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析 l 台词失败")
	}
	p.sendText(receiver, fmt.Sprintf("~%s~\n%s", resp.Data.Name, resp.Data.Content))
	return true, p.sendImage(receiver, resp.Data.Img)
}

func (p *DemosPlugin) handleSjkk(receiver *contact.Contact, arg string) (bool, error) {
	p.sendVoice(receiver, sjkkURL)
	return true, nil
}

// ==================== 搜索类处理器 ====================

func (p *DemosPlugin) handleBdbk(receiver *contact.Contact, arg string) (bool, error) {
	if p.Config.BdbkURL == "" {
		p.sendText(receiver, "百度百科功能未配置，请联系管理员配置 bdbk_url")
		return true, nil
	}
	if arg == "" {
		p.sendText(receiver, "请输入查询关键词，例如：百度百科 人工智能")
		return true, nil
	}
	body, err := p.httpGet(p.Config.BdbkURL + url.QueryEscape(arg))
	if err != nil {
		return true, err
	}
	var resp struct {
		Key        string `json:"key"`
		Desc       string `json:"desc"`
		Abstract   string `json:"abstract"`
		Copyrights string `json:"copyrights"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return true, fmt.Errorf("解析百科失败")
	}
	if resp.Key == "" {
		p.sendText(receiver, "没查到该百科含义")
		return true, nil
	}
	p.sendText(receiver, fmt.Sprintf("百科：%s\n%s\n摘要：%s\n版权：%s", resp.Key, resp.Desc, resp.Abstract, resp.Copyrights))
	return true, nil
}

func (p *DemosPlugin) handleSdj(receiver *contact.Contact, arg string) (bool, error) {
	if arg == "" {
		p.sendText(receiver, "请输入短剧名称，例如：搜短剧 爱情")
		return true, nil
	}
	body, err := p.httpGet(sdjURL + url.QueryEscape(arg))
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int `json:"code"`
		Data []struct {
			Name     string `json:"name"`
			Viewlink string `json:"viewlink"`
			Addtime  string `json:"addtime"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		p.sendText(receiver, "功能不可用了，请联系管理员")
		return true, nil
	}
	n := len(resp.Data)
	if n == 0 {
		p.sendText(receiver, "没有找到相关短剧")
		return true, nil
	}
	limit := p.Config.MaxList
	if n < limit {
		limit = n
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "查询到%d条\n", n)
	for i := 0; i < limit; i++ {
		item := resp.Data[i]
		fmt.Fprintf(&sb, "%d-> 短剧名称：【%s】\n短剧链接：%s\n添加时间：%s\n", i+1, item.Name, item.Viewlink, item.Addtime)
	}
	p.sendText(receiver, sb.String())
	return true, nil
}

// ==================== 生成类处理器 ====================

func (p *DemosPlugin) handleHunYanQingJian(receiver *contact.Contact, arg string) (bool, error) {
	defaultMsg := "婚宴请柬邀请函生成（仅供娱乐！）\n输入格式为（不需要【】符号）：婚宴请柬 【新郎名】,【新娘名】,【邀请人名】"
	if arg == "" {
		p.sendText(receiver, defaultMsg)
		return true, nil
	}
	parts := strings.Split(arg, ",")
	if len(parts) != 3 {
		p.sendText(receiver, defaultMsg)
		return true, nil
	}
	xl := strings.TrimSpace(parts[0])
	xn := strings.TrimSpace(parts[1])
	yq := strings.TrimSpace(parts[2])
	date := time.Now().AddDate(0, 0, 1).Format("20060102")
	imgURL := fmt.Sprintf("%s?to=%s&date=%s&head=%s|%s&event=结婚典礼恭备喜宴&time=中午十二时恭候&place=明月大酒店&tail=%s|%s",
		hunyanURL, url.QueryEscape(yq), date, url.QueryEscape(xl), url.QueryEscape(xn), url.QueryEscape(xl), url.QueryEscape(xn))
	return true, p.sendImage(receiver, imgURL)
}

func (p *DemosPlugin) handleDogDoc(receiver *contact.Contact, arg string) (bool, error) {
	defaultMsg := "狗屁不通文章生成\n输入格式为（不需要【】符号）：狗屁不通文章生成 +【内容描述】"
	if arg == "" {
		p.sendText(receiver, defaultMsg)
		return true, nil
	}
	body, err := p.httpGet(dogDocURL + url.QueryEscape(arg))
	if err != nil {
		return true, err
	}
	var resp struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		return true, fmt.Errorf("解析失败")
	}
	p.sendText(receiver, resp.Data)
	return true, nil
}

// ==================== 新增功能处理器 ====================

func (p *DemosPlugin) handleSjSing(receiver *contact.Contact, arg string) (bool, error) {
	body, err := p.httpGet(sjSingURL)
	if err != nil {
		return true, err
	}
	var resp struct {
		Code     int    `json:"code"`
		Nickname string `json:"nickname"`
		Pic      string `json:"pic"`
		Title    string `json:"title"`
		Singer   string `json:"singer"`
		Lyrics   string `json:"lyrics"`
		URL      string `json:"url"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Code != 200 {
		p.sendText(receiver, "获取随机唱失败")
		return true, nil
	}
	lines := strings.Split(resp.Lyrics, "\n")
	var sb strings.Builder
	sb.WriteString("\n[歌词]\n")
	for i, line := range lines {
		prefix := "> "
		if i%2 == 1 {
			prefix = "* "
		}
		fmt.Fprintf(&sb, "%s%s\n", prefix, line)
	}
	p.sendText(receiver, fmt.Sprintf("%s\n%s-%s\n%s", resp.Nickname, resp.Title, resp.Singer, sb.String()))
	if err := p.sendImage(receiver, resp.Pic); err != nil {
		slog.Warn("[demos] 随机唱头像发送失败", "err", err)
	}
	p.sendVoice(receiver, resp.URL)
	return true, nil
}

func (p *DemosPlugin) handleSearchMusic(receiver *contact.Contact, arg string) (bool, error) {
	if arg == "" {
		p.sendText(receiver, "输入格式为（不需要【】符号）：火子搜歌+【搜索词】")
		return true, nil
	}
	searchURL := kuwoSearchURL +
		"?prod=kwplayer_ar_9.3.7.2&corp=kuwo&newver=2&vipver=9.3.7.2" +
		"&source=kwplayer_ar_9.3.7.2_meizu.apk&p2p=1&notrace=0&client=kt" +
		"&pn=0&rn=45&ver=kwplayer_ar_9.3.7.2&vipver=1&show_copyright_off=1" +
		"&newver=2&correct=1&ft=music&cluster=0&strategy=2012&encoding=utf8" +
		"&rformat=json&vermerge=1&mobi=1&searchapi=2&issubtitle=1" +
		"&spPrivilege=0&all=" + url.QueryEscape(arg)

	headers := map[string]string{
		"User-Agent": "Dalvik/2.1.0 (Linux; U; Android 9; GM1910 Build/PQ3B.190801.07101020)",
		"Host":       "search.kuwo.cn",
	}
	body, err := p.httpGetWithHeaders(searchURL, headers)
	if err != nil {
		return true, err
	}
	var resp struct {
		Abslist []struct {
			SongName   string `json:"SONGNAME"`
			Artist     string `json:"ARTIST"`
			DCTargetID string `json:"DC_TARGETID"`
		} `json:"abslist"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil || len(resp.Abslist) == 0 {
		p.sendText(receiver, "没有找到相关歌曲")
		return true, nil
	}
	n := 10
	if len(resp.Abslist) < n {
		n = len(resp.Abslist)
	}
	var sb strings.Builder
	sb.WriteString("为你找到以下内容：\n")
	for i := 0; i < n; i++ {
		s := resp.Abslist[i]
		fmt.Fprintf(&sb, "%d. 【%s】\n%s\n点歌号：%s\n\n", i+1, s.SongName, s.Artist, s.DCTargetID)
	}
	p.sendText(receiver, sb.String())
	return true, nil
}

func (p *DemosPlugin) handleChooseMusic(receiver *contact.Contact, arg string) (bool, error) {
	if arg == "" {
		p.sendText(receiver, "输入格式为（不需要【】符号）：火子点歌+【点歌号】")
		return true, nil
	}
	chooseURL := kuwoMusicURL +
		"?f=web&source=kwplayerhd_ar_4.3.0.8_tianbao_T1A_qirui.apk" +
		"&type=convert_url_with_sign&br=128kmp3&rid=" + url.QueryEscape(arg)

	headers := map[string]string{
		"User-Agent":      "Apache-HttpClient/UNAVAILABLE (java 1.4)",
		"Accept-Encoding": "identity",
		"Host":            "nmobi.kuwo.cn",
	}
	body, err := p.httpGetWithHeaders(chooseURL, headers)
	if err != nil || body == "" {
		p.sendText(receiver, "点歌失败，也许点歌号不可用。")
		return true, nil
	}
	var musicResp struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &musicResp); err != nil || musicResp.Data.URL == "" {
		p.sendText(receiver, "点歌失败，也许点歌号不可用。")
		return true, nil
	}

	infoURL := kuwoInfoURL + "?httpsStatus=1&reqId=fcd6bc60-3e06-11ec-8722-67bb659a8433&musicId=" + url.QueryEscape(arg)
	infoHeaders := map[string]string{
		"User-Agent":      "Apache-HttpClient/UNAVAILABLE (java 1.4)",
		"Accept-Encoding": "identity",
		"Host":            "m.kuwo.cn",
	}
	infoBody, err := p.httpGetWithHeaders(infoURL, infoHeaders)
	if err != nil {
		p.sendText(receiver, "点歌失败，也许点歌号不可用。")
		return true, nil
	}
	var infoResp struct {
		Data struct {
			Songinfo struct {
				SongName string `json:"songName"`
				Artist   string `json:"artist"`
			} `json:"songinfo"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(infoBody), &infoResp); err != nil {
		p.sendText(receiver, "点歌失败，也许点歌号不可用。")
		return true, nil
	}
	p.sendText(receiver, fmt.Sprintf("%s - %s", infoResp.Data.Songinfo.SongName, infoResp.Data.Songinfo.Artist))
	p.sendVoice(receiver, musicResp.Data.URL)
	return true, nil
}

// ==================== 帮助处理器 ====================

func (p *DemosPlugin) handleExplain(receiver *contact.Contact, arg string) (bool, error) {
	msg := "粗略记录一下\n" +
		"【一言/今日句子】【小姐姐视频】【搜短剧】【百度百科】【撸猫】【旺财】【看星空】【温馨提示】【一句】【吟诗】【讲笑话】\n" +
		"【名画赏析】【婚宴请柬】【热点视频】【娱乐视频】【靓仔视频】【脑筋急转弯】【来段绕口令】\n" +
		"【谚语】【抽签】【答案之书】【吃什么】【king台词】【l台词】\n" +
		"【小黑子表情】【随机二次元】【保安日记】【acg美图】\n" +
		"【狗屁不通文章生成】【--小姐姐视频】【懒羊羊k歌】【怼脸自拍视频】【看穿搭】【丝滑舞蹈视频】【快手随机翻唱】【随机坤坤】\n" +
		"【随机唱】【火子搜歌】【火子点歌】"
	p.sendText(receiver, msg)
	return true, nil
}

func main() {
	p := &DemosPlugin{
		ConfigAbility: plugin.ConfigAbility[Config]{
			Config: Config{
				VideoNative: true,
				MaxList:     3,
			},
		},
		client: newHTTPClient(),
	}

	p.handlers = map[string]handlerFunc{
		// 图片
		"撸猫":    p.handleCat,
		"旺财":    p.handleDog,
		"看星空":   p.handleTw,
		"名画赏析":  p.handlePainting,
		"小黑子表情": p.handleXhzbq,
		"随机二次元": p.handleSjecy,
		"acg美图": p.handleAcg,

		// 视频
		"--小姐姐视频": p.handleXjj,
		"小姐姐视频":   p.handleXjj2,
		"热点视频":    p.handleRdVideo,
		"娱乐视频":    p.handleYlVideo,
		"靓仔视频":    p.handleBoyVideo,
		"懒羊羊k歌":   p.handleLyyKg,
		"怼脸自拍视频":  p.handleDuilian,
		"看穿搭":     p.handleChuanda,
		"丝滑舞蹈视频":  p.handleShwd,
		"快手随机翻唱":  p.handleKsfc,

		// 文本
		"温馨提示":   p.handleWxts,
		"一句":     p.handleYiju,
		"一言":     p.handleYiyan,
		"今日句子":   p.handleYiyan,
		"吟诗":     p.handleShici,
		"讲笑话":    p.handleHaha,
		"脑筋急转弯":  p.handleJzw,
		"来段绕口令":  p.handleRao,
		"谚语":     p.handleYanyu,
		"抽签":     p.handleChouqian,
		"答案之书":   p.handleBay,
		"吃什么":    p.handleEat,
		"保安日记":   p.handleRj,
		"king台词": p.handleKingTc,
		"l台词":    p.handleLolTc,
		"随机坤坤":   p.handleSjkk,

		// 搜索
		"百度百科": p.handleBdbk,
		"搜短剧":  p.handleSdj,

		// 生成
		"婚宴请柬":     p.handleHunYanQingJian,
		"狗屁不通文章生成": p.handleDogDoc,

		// 新增
		"随机唱":  p.handleSjSing,
		"火子搜歌": p.handleSearchMusic,
		"火子点歌": p.handleChooseMusic,

		// 帮助
		"demos": p.handleExplain,
	}

	slog.Info("[demos] 娱乐插件启动中...")
	plugin.Start(p)
}
