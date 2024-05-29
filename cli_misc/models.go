package cli_misc

type TSData struct {
	MediaId      string `json:"media_id"`
	ScriptPath   string `json:"script_path"`
	ScriptTSPath string `json:"script_ts_path"`
	LineTSPath   string `json:"line_ts_path"`
	VerseTSPath  string `json:"verse_ts_path"`
	Count        int    `json:"count"`
}
