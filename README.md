# Komodo Futurescapes Micro-Service Repository
A collection of web-service operations used to support various web/mobile applications

## Reverse Proxy Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getDesktopView </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getMobileView </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getResource </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - routeGET </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - routePOST </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - routePUT </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - routeDELETE </br>


## Authentication Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - login </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - logoff </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - validateOTP </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - validateEmail </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - validateNumber </br>


## User Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>
 
&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - createUser </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getInfo </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateInfo </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updatePreferences </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - deleteUser </br>


## Enhancement Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getUpgradeCollection </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getUpgradeDetails </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getReviews </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getRating </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - submitRating </br>


## Upgrade Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getProductCollection </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getProductDetails </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getReviews </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getRating </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - submitRating </br>


## Scheduling Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getAvailableTimes </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - scheduleService </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateScheduledService </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - cancelScheduledService </br>


## Invoice Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - createInvoice </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getInvoices </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateInvoice </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - deleteInvoice </br>


## Finance Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - submitPayment </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - schedulePayment </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - cancelPayment </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - subscribeReminders </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - unsubscribeReminders </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - refundPayment </br>


## Marketing Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - subscribeUser </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateMarketingPreferences </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - sendMarketingEmail </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - sendMarketingText </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - validateDeeplinkURL </br>


## News Media Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x  </br>
> Express Version: x.x.x </br>
> GraphQL Version: x.x.x [TBD][?] </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getArticle </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getAllArticles </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - postArticle </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - deleteArticle </br>


## Customer Support Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - submitFeedback </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getHelpTicket </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - createHelpTicket </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateHelpTicket </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - closeHelpTicket </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getHelpTopicList </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getHelpTopic </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - startBotChat </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - startCustSupptChat </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - sendChatMessage </br>


## Media Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getResource </br>


## WebApp Metrics Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - recordEvent </br>


## Hydroponic Distributed Control Service [Internal + External]
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x [TBD] </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - subscribe </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - unsubscribe </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getDefaultConfig </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getCustomConfig </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - setDeviceConfig </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getMinTemperature </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getMaxTemperature </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getMoistureLevel </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getLampBrightness </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getWaterSchedule </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getFeedSchedule </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getSleepSchedule </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportMechnicalFailure </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportSensorFailure </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportStasisFailure </br>


## Hydroponic Metrics Service [Internal]
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - recordEvent </br>


## Third-Party Services
> Azure Bot Service (Chat bots) </br>
> Azure Analytics Service (Logging/Analytics) </br>
> Azure IoT (Robotics/Automation) </br>
> MailChimp or ActiveCampaign (Marketing) [TEMP] </br>
> Freshbooks or Xero (Finances/Accounting/Invoices) [TEMP] </br>
