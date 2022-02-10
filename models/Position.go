package models

type PositionMap struct {
	Success       bool        `json:"success"`
	Msg           interface{} `json:"msg"`
	Code          int         `json:"code"`
	Content       Content     `json:"content"`
	Resubmittoken interface{} `json:"resubmitToken"`
	Requestid     interface{} `json:"requestId"`
}
type HrInfo struct {
	Userid       int         `json:"userId"`
	Portrait     string      `json:"portrait"`
	Realname     string      `json:"realName"`
	Positionname string      `json:"positionName"`
	Phone        interface{} `json:"phone"`
	Receiveemail interface{} `json:"receiveEmail"`
	Userlevel    string      `json:"userLevel"`
	Cantalk      bool        `json:"canTalk"`
}

type Positionnewlables struct {
}
type Result struct {
	Positionid            int               `json:"positionId"`
	Positionname          string            `json:"positionName"`
	Companyid             int               `json:"companyId"`
	Companyfullname       string            `json:"companyFullName"`
	Companyshortname      string            `json:"companyShortName"`
	Companylogo           string            `json:"companyLogo"`
	Companysize           string            `json:"companySize"`
	Industryfield         string            `json:"industryField"`
	Financestage          string            `json:"financeStage"`
	Companylabellist      []string          `json:"companyLabelList"`
	Firsttype             string            `json:"firstType"`
	Secondtype            string            `json:"secondType"`
	Thirdtype             string            `json:"thirdType"`
	Newfirsttype          interface{}       `json:"newFirstType"`
	Newsecondtype         interface{}       `json:"newSecondType"`
	Newthirdtype          interface{}       `json:"newThirdType"`
	Positionnewlables     Positionnewlables `json:"positionNewLables"`
	Skilllables           []interface{}     `json:"skillLables"`
	Positionlables        []interface{}     `json:"positionLables"`
	Industrylables        []interface{}     `json:"industryLables"`
	Createtime            string            `json:"createTime"`
	Formatcreatetime      string            `json:"formatCreateTime"`
	City                  string            `json:"city"`
	District              string            `json:"district"`
	Businesszones         interface{}       `json:"businessZones"`
	Salary                string            `json:"salary"`
	Salarymonth           string            `json:"salaryMonth"`
	Workyear              string            `json:"workYear"`
	Jobnature             string            `json:"jobNature"`
	Education             string            `json:"education"`
	Positionadvantage     string            `json:"positionAdvantage"`
	Imstate               string            `json:"imState"`
	Lastlogin             string            `json:"lastLogin"`
	Publisherid           int               `json:"publisherId"`
	Approve               int               `json:"approve"`
	Subwayline            string            `json:"subwayline"`
	Stationname           string            `json:"stationname"`
	Linestaion            string            `json:"linestaion"`
	Latitude              string            `json:"latitude"`
	Longitude             string            `json:"longitude"`
	Distance              interface{}       `json:"distance"`
	Hitags                interface{}       `json:"hitags"`
	Resumeprocessrate     int               `json:"resumeProcessRate"`
	Resumeprocessday      int               `json:"resumeProcessDay"`
	Score                 int               `json:"score"`
	Newscore              float64           `json:"newScore"`
	Matchscore            float64           `json:"matchScore"`
	Matchscoreexplain     interface{}       `json:"matchScoreExplain"`
	Query                 interface{}       `json:"query"`
	Explain               interface{}       `json:"explain"`
	Isschooljob           int               `json:"isSchoolJob"`
	Adword                int               `json:"adWord"`
	Plus                  interface{}       `json:"plus"`
	Pcshow                int               `json:"pcShow"`
	Appshow               int               `json:"appShow"`
	Deliver               int               `json:"deliver"`
	Gradedescription      interface{}       `json:"gradeDescription"`
	Promotionscoreexplain interface{}       `json:"promotionScoreExplain"`
	Ishothire             int               `json:"isHotHire"`
	Count                 int               `json:"count"`
	Aggregatepositionids  []interface{}     `json:"aggregatePositionIds"`
	Recalltype            interface{}       `json:"reCallType"`
	Userexpectid          int               `json:"userExpectId"`
	Userexpecttext        string            `json:"userExpectText"`
	Promotiontype         interface{}       `json:"promotionType"`
	Is51Job               bool              `json:"is51Job"`
	Expectjobid           int               `json:"expectJobId"`
	Encryptid             string            `json:"encryptId"`
	Positiondetail        string            `json:"positionDetail"`
	Positionaddress       string            `json:"positionAddress"`
	Hunterjob             bool              `json:"hunterJob"`
	Detailrecall          bool              `json:"detailRecall"`
	Famouscompany         bool              `json:"famousCompany"`
}
type Locationinfo struct {
	City                 interface{} `json:"city"`
	District             interface{} `json:"district"`
	Businesszone         interface{} `json:"businessZone"`
	Isallhotbusinesszone bool        `json:"isAllhotBusinessZone"`
	Locationcode         interface{} `json:"locationCode"`
	Querybygiscode       bool        `json:"queryByGisCode"`
}
type Queryanalysisinfo struct {
	Positionname  string      `json:"positionName"`
	Positionnames []string    `json:"positionNames"`
	Companyname   interface{} `json:"companyName"`
	Industryname  interface{} `json:"industryName"`
	Usefulcompany bool        `json:"usefulCompany"`
	Jobnature     interface{} `json:"jobNature"`
}
type Strategyproperty struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
type Categorytypeandname struct {
	Num3 string `json:"3"`
}
type Positionresult struct {
	Resultsize          int                 `json:"resultSize"`
	Result              []Result            `json:"result"`
	Locationinfo        Locationinfo        `json:"locationInfo"`
	Queryanalysisinfo   Queryanalysisinfo   `json:"queryAnalysisInfo"`
	Strategyproperty    Strategyproperty    `json:"strategyProperty"`
	Hotlabels           interface{}         `json:"hotLabels"`
	Hitags              interface{}         `json:"hiTags"`
	Benefittags         interface{}         `json:"benefitTags"`
	Industryfield       interface{}         `json:"industryField"`
	Companysize         interface{}         `json:"companySize"`
	Positionname        interface{}         `json:"positionName"`
	Totalcount          int                 `json:"totalCount"`
	Totalscore          float64             `json:"totalScore"`
	Triggerorsearch     bool                `json:"triggerOrSearch"`
	Categorytypeandname Categorytypeandname `json:"categoryTypeAndName"`
	Categorytagcodes    []interface{}       `json:"categoryTagCodes"`
	Tagcodes            []interface{}       `json:"tagCodes"`
}
type Content struct {
	Showid         string         `json:"showId"`
	Hrinfomap      map[string]HrInfo      `json:"hrInfoMap"`
	Pageno         int            `json:"pageNo"`
	Positionresult Positionresult `json:"positionResult"`
	Pagesize       int            `json:"pageSize"`
}