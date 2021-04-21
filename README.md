# Komodo Futurescapes Microservice Repository
A collection of microservices used to support various web/mobile applications

## Gateway Router Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/gateway-router/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /route </br>


## Gateway Caching Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/cache/gateway-cache/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /desktop </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: /desktop/{key} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /desktop/{key} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /mobile </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: /mobile/{key} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /mobile/{key} </br>


## Web Security Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/security/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /auth/login </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /auth/logoff </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /validate/otp </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /validate/email </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /validate/phone </br>


## Catalog Search Service
> Release Version: 0.1 (beta) </br>
> Golang Version: 1.16 </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com:443/catalog-search/v0.1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /search?category={categoryID}&keyword={keyword} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /category/{id} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /item/{id} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /item/{itemID}/review/{reviewID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /item/{itemID}/reviews </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /item/{itemID}/submit/review </br>


## Catalog Caching Service
> Release Version: 0.1 (beta) </br>
> Golang Version: 1.16 </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com:443/cache/catalog-cache/v0.1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /{key} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: /{key} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /{key} </br>


## User Management Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/user-mgmt/v1/...
 
&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /user/create </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /user/details </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:    /users?firstName={firstName}&lastName={lastName}&role={roleType}&geoloc={region} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /user/update </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /user/delete </br>


## Order Details Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/order-details/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /create </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /order/{id}/details </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - PUT:    /order/{id}//update </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /order/{id}/delete </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /history </br>


## Order Scheduling Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/order-scheduling/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: schedule/delivery </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: schedule/arrival </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: schedule/availibility </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: schedule/update </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: schedule/cancel </br>


## Customer Payment Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/cust-payment/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /submit </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /schedule </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /cancel </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /refund </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /cancel </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /reminders </br>


## Customer Marketing Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/cust-marketing/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /subscribe </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /preferences </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /post/email </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /post/text </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /validate/deeplink </br>


## Customer Outreach (News) Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x  </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/cust-relations/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:    /article/{articleID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:    /articles?id={catID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /article </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /article </br>


## Customer Support Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/cust-support/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /ticket/list </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /ticket/details </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /ticket/create </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /ticket/attach </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /ticket/close </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /ticket/submit </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:    /topic/list?id={listID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:    /topic/{topicID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /chat/start </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /chat/post </br>


## Customer Support Caching Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/cache/cust-support-cache/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - PUT: / </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: /{key} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /{key} </br>


## HR Careers Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/hr-careers/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: TODO </br>


## HR Job Listing Search Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/hr-listing-search/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: /listing/{id} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: /listing/keyword?={keyword}&dept={dept}&loc={loc}&type={type}&date={date}&exp={exp}&remote={remote} </br>


## Serverside Rendering Service (Rendering Engine)
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/render-engine/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /resource </br>


## Web Asset Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/assets/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: resource/{type}/{ID} </br>


## Web Configuration Service
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/config/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /app </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /user </br>


## Web Metrics API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/metrics/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /record </br>


## Cache Cleaner Utility Service
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x ??? </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://api.{env}.somecloudurl.com/util/cache-cleaner/v1/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /start </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /stop </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /execute </br>


## Third-Party Services
> Azure Bot Service (Chat bots) </br>
> Azure Analytics Service (Logging/Analytics) </br>
> Some sort of resume analysis tool
> Azure IoT (Robotics/Automation) </br>
> MailChimp or ActiveCampaign (Marketing) </br>
> Freshbooks or Xero (Finances/Accounting/Invoices) </br>


## Resources
> https://lucid.app/lucidchart/invitations/accept/4f240be9-e7e3-47e5-8d0e-9044688a704a 