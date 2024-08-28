package osu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RankingType string

var (
	performanceRankValue RankingType = "performance"
	chartsRankValue      RankingType = "charts"
	countryRankValue     RankingType = "country"
	scoreRankValue       RankingType = "score"

	PerformanceRank *RankingType = &performanceRankValue
	ChartsRank      *RankingType = &chartsRankValue
	CountryRank     *RankingType = &countryRankValue
	ScoreRank       *RankingType = &scoreRankValue
)

func (client *OsuAPIClient) GetRanking(modeStringInput *ModeString, rankingTypeInput *RankingType) (*RankedReturn, error) {
	// {{baseUrl}}/api/v2/rankings/:mode/:type?filter=all

	// if modeString is nil we default to osuStandard since thats what people usually play
	var modeString string
	if modeStringInput != nil {
		modeString = string(*modeStringInput)
	} else {
		modeString = "osu"
	}

	//same pattern as above but for rankingType
	var rankingType string
	if rankingTypeInput != nil {
		rankingType = string(*rankingTypeInput)
	} else {
		rankingType = "performance"
	}

	var url = fmt.Sprintf("https://osu.ppy.sh/api/v2/rankings/%s/%s", modeString, rankingType)

	request, err := http.NewRequest("get", url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %s", err.Error())
	}
	ccTokenDeref := client.ccToken
	//let's check that ccToken actually has a value before using it.
	if ccTokenDeref == "" {
		return nil, errors.New("failed to get cc token")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+ccTokenDeref)
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	var responseData RankedReturn
	responseDataBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %s", err.Error())
	}
	err = json.Unmarshal(responseDataBytes, &responseData)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %s", err.Error())
	}
	return &responseData, nil
}

type RankedReturn struct {
	Cursor struct {
		Page int `json:"page"`
	} `json:"cursor"`
	Ranking []struct {
		Count100  int `json:"count_100"`
		Count300  int `json:"count_300"`
		Count50   int `json:"count_50"`
		CountMiss int `json:"count_miss"`
		Level     struct {
			Current  int `json:"current"`
			Progress int `json:"progress"`
		} `json:"level"`
		GlobalRank             int     `json:"global_rank"`
		GlobalRankExp          any     `json:"global_rank_exp"`
		Pp                     float64 `json:"pp"`
		PpExp                  int     `json:"pp_exp"`
		RankedScore            int64   `json:"ranked_score"`
		HitAccuracy            float64 `json:"hit_accuracy"`
		PlayCount              int     `json:"play_count"`
		PlayTime               int     `json:"play_time"`
		TotalScore             int64   `json:"total_score"`
		TotalHits              int     `json:"total_hits"`
		MaximumCombo           int     `json:"maximum_combo"`
		ReplaysWatchedByOthers int     `json:"replays_watched_by_others"`
		IsRanked               bool    `json:"is_ranked"`
		GradeCounts            struct {
			Ss  int `json:"ss"`
			SSH int `json:"ssh"`
			S   int `json:"s"`
			Sh  int `json:"sh"`
			A   int `json:"a"`
		} `json:"grade_counts"`
		RankChangeSince30Days int `json:"rank_change_since_30_days"`
		User                  struct {
			AvatarURL     string    `json:"avatar_url"`
			CountryCode   string    `json:"country_code"`
			DefaultGroup  string    `json:"default_group"`
			ID            int       `json:"id"`
			IsActive      bool      `json:"is_active"`
			IsBot         bool      `json:"is_bot"`
			IsDeleted     bool      `json:"is_deleted"`
			IsOnline      bool      `json:"is_online"`
			IsSupporter   bool      `json:"is_supporter"`
			LastVisit     time.Time `json:"last_visit"`
			PmFriendsOnly bool      `json:"pm_friends_only"`
			ProfileColour any       `json:"profile_colour"`
			Username      string    `json:"username"`
			Country       struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"country"`
			Cover struct {
				CustomURL string `json:"custom_url"`
				URL       string `json:"url"`
				ID        string `json:"id"`
			} `json:"cover"`
		} `json:"user"`
	} `json:"ranking"`
	Total int `json:"total"`
}
