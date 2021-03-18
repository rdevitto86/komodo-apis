# Komodo Futurescapes Micro-Service Repository
A collection of web-service operations used to support various web/mobile applications

## Gateway Router API**
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/gateway-router/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /route </br>


## Gateway Caching API**
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/gateway-cache/{version}/cache...

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /desktop </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /mobile </br>


## Web Security API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/web-security-api/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /auth/login </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /auth/logoff </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /validate/otp </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /validate/email </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /validate/phone </br>


## Catalog Item API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/catalog-api/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /product/{catalogID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /service/{catalogID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /category?id={catID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /ratings/{catalogID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /reviews/{catalogID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /submitRating </br>


## Catalog Search API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/catalog-search-api/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:  /search?category={category}&keyword={keyword} </br>


## Catalog Caching API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/cust-caching-api/{version}/cache/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: /item/{key} </br>


## User Details API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/user-details-api/{version}/user/...
 
&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /create </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /details </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /update </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /delete </br>


## Order Details API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/order-details-api/{version}/order/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /create </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /details </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /list </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /update </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /delete </br>


## Order Scheduling API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/order-scheduling-api/{version}/schedule/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /delivery </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /arrival </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /availibility </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /update </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /cancel </br>


## Finance API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/finance-api/{version}/payment/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /submit </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /schedule </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /cancel </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /refund </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /cancel </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /reminders </br>


## Marketing API
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/marketing-api/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /subscribe </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /preferences </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /post/email </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /post/text </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /validate/deeplink </br>


## Customer Relations API
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x  </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/cust-relations-api/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:    /article/{articleID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET:    /articles?id={catID} </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST:   /article </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - DELETE: /article </br>


## Customer Support API
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/cust-support-api/{version}/...

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


## Customer Support Caching API
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/cust-support-cache-api/{version}/cache/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: /topic/{key} </br>


## Web Resource API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/web-resource-api/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /resource </br>


## Web Asset API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/web-asset-api/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - GET: resource/{type}/{ID} </br>


## Web Configuration API
> Release Version: 0.1 (alpha) </br>
> Golang Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/web-config-api/{version}/config/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /app </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /user </br>


## Web Metrics API
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/web-metrics-api/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /record </br>


## Cache Cleaner Utility API
> Release Version: 0.1 (alpha) </br>
> Node Version: x.x.x </br>
> Express Version: x.x.x </br>

&nbsp;&nbsp;&nbsp; **Host Domain:**
> https://someazurelink/util/cache-cleaner/{version}/...

&nbsp;&nbsp;&nbsp; **API Operations:** </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /start </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /stop </br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  - POST: /execute </br>


## Third-Party Services
> Azure Bot Service (Chat bots) </br>
> Azure Analytics Service (Logging/Analytics) </br>
> Azure IoT (Robotics/Automation) </br>
> MailChimp or ActiveCampaign (Marketing) [TEMP] </br>
> Freshbooks or Xero (Finances/Accounting/Invoices) [TEMP] </br>
