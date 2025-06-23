package speedtester

type SpeedTestResult struct {
	UploadSpeedMbps   float64 `json:"upload_speed_mbps"`  
	DownloadSpeedMbps float64 `json:"download_speed_mbps"` 
}