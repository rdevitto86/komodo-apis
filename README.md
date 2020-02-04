# Komodo Future Solutions Web-Service Repository
A collection of web-service operations used to support various web/mobile applications

## WebApp Delegator Service 
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - serveECOMDesktop </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - serveECOMMobile </br>


## Authentication Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - login </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - logoff </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - validateOTP </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - validateEmail </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - validateMobile </br>


## User Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x  </br>
> [?] GraphQL Version: x.x.x </br>
 
&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - createUser </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getUserInfo </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateUserInfo </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateUserPreferences </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - deleteUser </br>


## Merchandise Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x  </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getProductCollection </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getProductDetails </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getReviews </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getRating </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - submitRating </br>


## ContractedServices Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getServiceCollection </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getServiceDetails </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getReviews </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getRating </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - submitRating </br>


## Order Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x  </br>
> [?] GraphQL Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - createOrder </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getOrderList </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateOrder </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - cancelOrder </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - scheduleService </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateScheduledService </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - cancelScheduledService </br>


## Finance Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x  </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - submitPayment </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getInvoiceDetails </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updatePaymentInfo </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - sendPaymentReminder </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - refundPayment </br>


## Marketing Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - updateUserMarketingPreferences </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - unsubscribeUser </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - sendMarketingEmail </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - sendMarketingText </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - validateDeeplinkURL </br>


## News (Blog) Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x  </br>
> Express Version: x.x.x </br>
> [?] GraphQL Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getArticle </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getAllArticles </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - postArticle </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - deleteArticle </br>


## Customer-Support Service
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


## Hydroponic Support Unit (HSU) Service [Internal]
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x  </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - subscribeHSU </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - unSubscribeHSU </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - getLatestHSUConfig </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportTemperature </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportMoisture </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportLampBrightness </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportWateringTime </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportFeedingTime </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - reportMechnicalFailure </br>


## Hydroponics Administration Service [Internal]
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x  </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - overrideTemperature </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - overrideMoisture </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - overrideLampBrightness </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - overrideWateringSchedule </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - overrideFeedingSchedule </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - overrideSleepSchedule </br>


## Internal Operations Analytics Service [Internal]
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - recordEvent </br>


## WebApp Analytics Service
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
