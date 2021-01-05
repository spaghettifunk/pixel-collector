package models

import "github.com/google/uuid"

// PayloadRaw defines the object in which the data is collected from the website
type PayloadRaw struct {
	TrackingID       string `json:"tracking_id" query:"id" description:"id for the app/website you are tracking"`
	UserID           string `json:"user_id" query:"uid" description:"id of the user"`
	Timestamp        string `json:"timestamp" query:"ts" description:"timestamp of the incoming message"`
	Type             string `json:"type" query:"ev" description:"the event that is being triggered"`
	Version          string `json:"version" query:"v" description:"openpixel js version number"`
	DocumentLocation string `json:"document_loc" query:"dl" description:"document location"`
	ReferrerLocation string `json:"referrer_loc" query:"rl" description:"referrer location"`
	DocumentEncoding string `json:"document_enc" query:"de" description:"document encoding"`
	ScreenResolution string `json:"screen_res" query:"sr" description:"screen resolution"`
	Viewport         string `json:"viewport" query:"vp" description:"viewport"`
	ColorDepth       string `json:"color_depth" query:"cd" description:"color depth"`
	DocumentTitle    string `json:"document_title" query:"dt" description:"document title"`
	BrowserName      string `json:"browser_name" query:"bn" description:"browser name"`
	MobileDevice     string `json:"mobile_device" query:"md" description:"mobile device"`
	UserAgent        string `json:"user_agent" query:"ua" description:"full user agent"`
	TimeZone         string `json:"timezone" query:"tz" description:"timezone offset (minutes away from utc)"`
	UTMSource        string `json:"utm_source" query:"utm_source" description:"campaign source"`
	UTMMedium        string `json:"utm_medium" query:"utm_medium" description:"campaign medium"`
	UTMTerm          string `json:"utm_term" query:"utm_term" description:"campaign Term"`
	UTMContent       string `json:"utm_content" query:"utm_content" description:"campaign content"`
	UTMCampaign      string `json:"utm_campaign" query:"utm_campaign" description:"campaign name"`
}

// Payload is the object sent to the bus
type Payload struct {
	ID uuid.UUID `json:"id" description:"uuid of the single event"`
	PayloadRaw
}
