#Cookie

|params|ログイン|備考|
|----|----|----|
|__utma|✗|Google Analytics|
|__utmb|✗|Google Analytics|
|__utmc|✗|Google Analytics|
|__utmz|✗|Google Analytics|
|sessionid|✗|そのまんま|
|steamCC_182_167_31_201|✗|終盤はIP 内容はCCコード (JP)|
|timezoneOffset|✗|そのまんま|
|570_17workshopQueueTime|✓|1420275955|
|Steam_Language|✓|japanese|
|recentlyVisitedAppHub|✓|複数のAppIDをコンマ区切りにし、URLエンコードしたのが入ってる|
|steamLogin|✓|SteamID64 + "&#124;&#124;" + 42桁のHex|
|webTradeEligibility|✓|jsonをURLエンコードしてたのが入ってる|

##webTradeEligibility
```
{
	"allowed": 0,
	"reason": 0,
	"allowed_at_time": 0,
	"steamguard_required_days": 0,
	"sales_this_year": 0,
	"max_sales_per_year": 0,
	"forms_requested": 0,
	"new_device_cooldown_days": 0,
	"time_payments_trusted": 0
}
```

多分重要なのは`steamLogin`と`sessionid`
