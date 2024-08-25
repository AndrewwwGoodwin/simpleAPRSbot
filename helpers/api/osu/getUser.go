package osu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (client *OsuAPIClient) GetUser(user string, mode *ModeString, key *GetUserKeyType) (*GetUserReturnData, error) {
	err := client.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %v", err)
	}
	var modeString string
	switch mode {
	case nil:
		modeString = ""
	default:
		modeString = string(*mode)
	}

	var keyString string
	switch key {
	case nil:
		keyString = ""
	default:
		keyString = string(*key)
	}

	// if mode is specified, lets use it.
	if keyString == "username" {
		if strings.HasPrefix(user, "@") {
			user = "@" + user
		}
	} else if keyString == "id" {
		if _, err := strconv.Atoi(user); err != nil {
			return nil, errors.New("invalid user id")
		}
	}

	//begin building the request here
	var requestUrl = fmt.Sprintf("https://osu.ppy.sh/api/v2/users/%s/%s", user, modeString)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	ccTokenDeref := client.ccToken
	//let's check that ccToken actually has a value before using it.
	if ccTokenDeref == "" {
		return nil, errors.New("failed to get cc token")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+ccTokenDeref)

	query := request.URL.Query()
	query.Add("key", keyString)
	request.URL.RawQuery = query.Encode()

	// now that the request is built, lets send it!
	requestData, err := client.httpClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(requestData.Body)
	// read the information from the request
	data, err := io.ReadAll(requestData.Body)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	var returnData = GetUserReturnData{}
	// unmarshal that into GetUserReturnData
	err = json.Unmarshal(data, &returnData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}
	// and finally send it back to wherever asked for it.
	return &returnData, nil
}

type GetUserKeyType string

var (
	keyIdValue       GetUserKeyType = "id"
	keyUsernameValue GetUserKeyType = "username"
	KeyID                           = &keyIdValue
	KeyUsername                     = &keyUsernameValue
)

type GetUserReturnData struct {
	AvatarURL     string      `json:"avatar_url"`
	CountryCode   string      `json:"country_code"`
	DefaultGroup  string      `json:"default_group"`
	ID            int         `json:"id"`
	IsActive      bool        `json:"is_active"`
	IsBot         bool        `json:"is_bot"`
	IsDeleted     bool        `json:"is_deleted"`
	IsOnline      bool        `json:"is_online"`
	IsSupporter   bool        `json:"is_supporter"`
	LastVisit     interface{} `json:"last_visit"`
	PmFriendsOnly bool        `json:"pm_friends_only"`
	ProfileColour interface{} `json:"profile_colour"`
	Username      string      `json:"username"`
	CoverURL      string      `json:"cover_url"`
	Discord       string      `json:"discord"`
	HasSupported  bool        `json:"has_supported"`
	Interests     interface{} `json:"interests"`
	JoinDate      time.Time   `json:"join_date"`
	Location      string      `json:"location"`
	MaxBlocks     int         `json:"max_blocks"`
	MaxFriends    int         `json:"max_friends"`
	Occupation    interface{} `json:"occupation"`
	Playmode      string      `json:"playmode"`
	Playstyle     []string    `json:"playstyle"`
	PostCount     int         `json:"post_count"`
	ProfileHue    int         `json:"profile_hue"`
	ProfileOrder  []string    `json:"profile_order"`
	Title         interface{} `json:"title"`
	TitleURL      interface{} `json:"title_url"`
	Twitter       string      `json:"twitter"`
	Website       string      `json:"website"`
	Country       struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"country"`
	Cover struct {
		CustomURL string `json:"custom_url"`
		URL       string `json:"url"`
		ID        string `json:"id"`
	} `json:"cover"`
	Kudosu struct {
		Available int `json:"available"`
		Total     int `json:"total"`
	} `json:"kudosu"`
	AccountHistory          []interface{} `json:"account_history"`
	ActiveTournamentBanner  interface{}   `json:"active_tournament_banner"`
	ActiveTournamentBanners []interface{} `json:"active_tournament_banners"`
	Badges                  []interface{} `json:"badges"`
	BeatmapPlaycountsCount  int           `json:"beatmap_playcounts_count"`
	CommentsCount           int           `json:"comments_count"`
	DailyChallengeUserStats struct {
		DailyStreakBest     int       `json:"daily_streak_best"`
		DailyStreakCurrent  int       `json:"daily_streak_current"`
		LastUpdate          time.Time `json:"last_update"`
		LastWeeklyStreak    time.Time `json:"last_weekly_streak"`
		Playcount           int       `json:"playcount"`
		Top10PPlacements    int       `json:"top_10p_placements"`
		Top50PPlacements    int       `json:"top_50p_placements"`
		UserID              int       `json:"user_id"`
		WeeklyStreakBest    int       `json:"weekly_streak_best"`
		WeeklyStreakCurrent int       `json:"weekly_streak_current"`
	} `json:"daily_challenge_user_stats"`
	FavouriteBeatmapsetCount int           `json:"favourite_beatmapset_count"`
	FollowerCount            int           `json:"follower_count"`
	GraveyardBeatmapsetCount int           `json:"graveyard_beatmapset_count"`
	Groups                   []interface{} `json:"groups"`
	GuestBeatmapsetCount     int           `json:"guest_beatmapset_count"`
	LovedBeatmapsetCount     int           `json:"loved_beatmapset_count"`
	MappingFollowerCount     int           `json:"mapping_follower_count"`
	MonthlyPlaycounts        []struct {
		StartDate string `json:"start_date"`
		Count     int    `json:"count"`
	} `json:"monthly_playcounts"`
	NominatedBeatmapsetCount int `json:"nominated_beatmapset_count"`
	Page                     struct {
		HTML string `json:"html"`
		Raw  string `json:"raw"`
	} `json:"page"`
	PendingBeatmapsetCount int      `json:"pending_beatmapset_count"`
	PreviousUsernames      []string `json:"previous_usernames"`
	RankHighest            struct {
		Rank      int       `json:"rank"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"rank_highest"`
	RankedBeatmapsetCount int `json:"ranked_beatmapset_count"`
	ReplaysWatchedCounts  []struct {
		StartDate string `json:"start_date"`
		Count     int    `json:"count"`
	} `json:"replays_watched_counts"`
	ScoresBestCount   int `json:"scores_best_count"`
	ScoresFirstCount  int `json:"scores_first_count"`
	ScoresPinnedCount int `json:"scores_pinned_count"`
	ScoresRecentCount int `json:"scores_recent_count"`
	Statistics        struct {
		Count100  int `json:"count_100"`
		Count300  int `json:"count_300"`
		Count50   int `json:"count_50"`
		CountMiss int `json:"count_miss"`
		Level     struct {
			Current  int `json:"current"`
			Progress int `json:"progress"`
		} `json:"level"`
		GlobalRank             int         `json:"global_rank"`
		GlobalRankExp          interface{} `json:"global_rank_exp"`
		Pp                     float64     `json:"pp"`
		PpExp                  int         `json:"pp_exp"`
		RankedScore            int64       `json:"ranked_score"`
		HitAccuracy            float64     `json:"hit_accuracy"`
		PlayCount              int         `json:"play_count"`
		PlayTime               int         `json:"play_time"`
		TotalScore             int64       `json:"total_score"`
		TotalHits              int         `json:"total_hits"`
		MaximumCombo           int         `json:"maximum_combo"`
		ReplaysWatchedByOthers int         `json:"replays_watched_by_others"`
		IsRanked               bool        `json:"is_ranked"`
		GradeCounts            struct {
			Ss  int `json:"ss"`
			SSH int `json:"ssh"`
			S   int `json:"s"`
			Sh  int `json:"sh"`
			A   int `json:"a"`
		} `json:"grade_counts"`
		CountryRank int `json:"country_rank"`
		Rank        struct {
			Country int `json:"country"`
		} `json:"rank"`
	} `json:"statistics"`
	SupportLevel     int `json:"support_level"`
	UserAchievements []struct {
		AchievedAt    time.Time `json:"achieved_at"`
		AchievementID int       `json:"achievement_id"`
	} `json:"user_achievements"`
	RankHistory struct {
		Mode string `json:"mode"`
		Data []int  `json:"data"`
	} `json:"rank_history"`
	RankHistory2 struct {
		Mode string `json:"mode"`
		Data []int  `json:"data"`
	} `json:"rankHistory"`
	RankedAndApprovedBeatmapsetCount int `json:"ranked_and_approved_beatmapset_count"`
	UnrankedBeatmapsetCount          int `json:"unranked_beatmapset_count"`
}
