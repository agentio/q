# q

`q` is a command line tool for working with Google's
[Service Infrastructure](https://cloud.google.com/service-infrastructure/docs/overview) APIs.

To install `q` on any system with Go installed, run the following:
```
go install git@github.com/agentio/q:latest
```

## quickstart

Here's an easy way to explore the Service Infrastructure APIs.
We'll use `q` in the Google Cloud Shell to set up a managed service
on Google Cloud Run.

This assumes that you have a Google Cloud account with a project
created and access to the Google Cloud Shell. You don't have to
start with a fresh project, but that would be a good way to keep
everything associated with your demo in one easy-to-find place.

After opening Google Cloud Shell, install `q` with:

`go install github.com/agentio/q@latest`

When that finishes, run `q` with no arguments.

```
$ q
Manage APIs with Service Infrastructure

Usage:
  q [command]

Available Commands:
  api-keys           Manage API keys with the API Keys API
  compile            Compile a Service Configuration for an API
  completion         Generate the autocompletion script for the specified shell
  demo               Set up a sample managed service
  doctor             Verify necessary dependencies and configuration
  help               Help about any command
  inspect            Read API information from compiled file descriptors
  jwt                Read, verify, and generate JSON Web Tokens
  logging            Write and manage log entries with the Cloud Logging API
  monitoring         Monitor services with the Cloud Monitoring API
  service-control    Control API services with the Service Control API
  service-management Manage service descriptions with the Service Management API
  service-usage      Manage usage of APIs with the Service Usage API
  translate          Translate with the Cloud Translation API

Flags:
  -h, --help   help for q

Use "q [command] --help" for more information about a command.
```

Now run `q doctor` to check your `gcloud` configuration.
You should see something like this:

```
$ q doctor
Checking gcloud configuration...
account = YOUR_EMAIL
project = YOUR_PROJECT
Error: run/region required
```

If you need to set a run region,
do it with the following command,
and if you don't have a preferred
region, use `us-west1`.

```
$ gcloud config set run/region us-west1
Updated property [run/region].
```

Now rerun `q doctor`.

```
$ q doctor
Checking gcloud configuration...
account = YOUR_EMAIL
project = YOUR_PROJECT
run/region = us-west1
Checking application default credentials...
token = ya29.[REDACTED]
Everything looks good!
```

Next run `q demo`. This generates a directory of
files that you'll use to set up your service.
They will be customized to your `gcloud` project
and run region, so you won't need to edit them.

```
$ q demo
To run the demo, see the SETUP.sh script in stores-demo
```

Go inside the generated directory and look around.
```
$ cd stores-demo
```

Run the `SETUP.sh` script to configure your project and create your instance.
```
$ sh SETUP.sh 
Enabling the Service Control API.
Enabling the Service Management API.
Enabling the Cloud Run Admin API.
Operation "operations/acf.p2-622059496897-baab248d-0eae-4d5e-9e65-4fe1aa26bae0" finished successfully.
Creating the service from the service descriptor and API config file.
Waiting for async operation operations/serviceConfigs.stores.endpoints.YOUR_PROJECT.cloud.goog:7bbc869f-df3b-40fd-88f7-26bfe4c0c053 to complete...
Operation finished successfully. The following command can describe the Operation details:
 gcloud endpoints operations describe operations/serviceConfigs.stores.endpoints.YOUR_PROJECT.cloud.goog:7bbc869f-df3b-40fd-88f7-26bfe4c0c053

Waiting for async operation operations/rollouts.stores.endpoints.YOUR_PROJECT.cloud.goog:8340c75f-411a-4c5c-a282-1c1312cee543 to complete...
Operation finished successfully. The following command can describe the Operation details:
 gcloud endpoints operations describe operations/rollouts.stores.endpoints.YOUR_PROJECT.cloud.goog:8340c75f-411a-4c5c-a282-1c1312cee543

Service Configuration [2024-10-29r1] uploaded for service [stores.endpoints.YOUR_PROJECT.cloud.goog]

To manage your API, go to: https://console.cloud.google.com/endpoints/api/stores.endpoints.YOUR_PROJECT.cloud.goog/overview?project=YOUR_PROJECT
Creating a service account to run the server container.
ERROR: (gcloud.iam.service-accounts.create) Resource in projects [YOUR_PROJECT] is the subject of a conflict: Service account stores already exists within project projects/YOUR_PROJECT.
- '@type': type.googleapis.com/google.rpc.ResourceInfo
  resourceName: projects/YOUR_PROJECT/serviceAccounts/stores@YOUR_PROJECT.iam.gserviceaccount.com
Giving the service account roles to call the Service Control APIs.
Updated IAM policy for project [YOUR_PROJECT].
bindings:
- members:
  - serviceAccount:stores@YOUR_PROJECT.iam.gserviceaccount.com
  role: roles/cloudtrace.agent
- members:
  - serviceAccount:service-622059496897@containerregistry.iam.gserviceaccount.com
  role: roles/containerregistry.ServiceAgent
- members:
  - serviceAccount:622059496897-compute@developer.gserviceaccount.com
  role: roles/editor
- members:
  - user:YOUR_EMAIL
  role: roles/owner
- members:
  - serviceAccount:service-622059496897@gcp-sa-pubsub.iam.gserviceaccount.com
  role: roles/pubsub.serviceAgent
- members:
  - serviceAccount:service-622059496897@serverless-robot-prod.iam.gserviceaccount.com
  role: roles/run.serviceAgent
- members:
  - serviceAccount:stores@YOUR_PROJECT.iam.gserviceaccount.com
  role: roles/servicemanagement.serviceController
etag: BwYlom1WNEs=
version: 1
Updated IAM policy for project [YOUR_PROJECT].
bindings:
- members:
  - serviceAccount:stores@YOUR_PROJECT.iam.gserviceaccount.com
  role: roles/cloudtrace.agent
- members:
  - serviceAccount:service-622059496897@containerregistry.iam.gserviceaccount.com
  role: roles/containerregistry.ServiceAgent
- members:
  - serviceAccount:622059496897-compute@developer.gserviceaccount.com
  role: roles/editor
- members:
  - user:YOUR_EMAIL
  role: roles/owner
- members:
  - serviceAccount:service-622059496897@gcp-sa-pubsub.iam.gserviceaccount.com
  role: roles/pubsub.serviceAgent
- members:
  - serviceAccount:service-622059496897@serverless-robot-prod.iam.gserviceaccount.com
  role: roles/run.serviceAgent
- members:
  - serviceAccount:stores@YOUR_PROJECT.iam.gserviceaccount.com
  role: roles/servicemanagement.serviceController
etag: BwYlom1zMec=
version: 1
Creating the Cloud Run container.
Applying new configuration to Cloud Run service [stores] in project [YOUR_PROJECT] region [us-west1]
OK Deploying new service... Done.                                                                                                                                                                                                                                
  OK Creating Revision...                                                                                                                                                                                                                                        
  OK Routing traffic...                                                                                                                                                                                                                                          
Done.                                                                                                                                                                                                                                                            
New configuration has been applied to service [stores].
URL: https://YOUR_HOST.us-west1.run.app
Configuring IAM to allow outside access to the container.
Updated IAM policy for service [stores].
bindings:
- members:
  - allUsers
  role: roles/run.invoker
etag: BwYlonNoIzE=
version: 1
```

Now go to the Cloud Run section of the Google Cloud Console, select your service, and find the service URL.
Click on that to call your service from your browser. You'll get a response like this:

```
{"code":404,"message":"The current request is not defined by this API."}
```

That's good, because it's true. The root path is not defined by your API. To call your API, add `/v1/stores` to the path.

You'll get JSON like this:
```
{"stores":[{"name":"stores/0","type":"office","title":"Columbus, NM 88029","location":{"latitude":31.8301201,"longitude":-107.638199},"address":{"street":"South Main Street","regionCode":"us"}},{"name":"stores/1","type":"office","title":"United States Post Office - Gateway Station","location":{"latitude":37.7949409,"longitude":-122.399437},"address":{"zipCode":94111,"regionCode":"us"}},{"name":"stores/2","type":"office","title":"Belington Post Office","location":{"latitude":39.0219955,"longitude":-79.9369202},"address":{"regionCode":"us"}},{"name":"stores/3","type":"office","title":"Denali National Park Post Office","location":{"latitude":63.7302589,"longitude":-148.892105},"address":{"regionCode":"us"}},{"name":"stores/4","type":"office","title":"Noti Post Office","location":{"latitude":44.057888,"longitude":-123.448227},"address":{"street":"Highway 126","city":"Noti","state":"OR","regionCode":"us"}},{"name":"stores/5","type":"office","title":"Piedra Post Office","location":{"latitude":36.8229141,"longitude":-119.364891},"address":{"regionCode":"us"}},{"name":"stores/6","type":"office","title":"Lafayette Post Office","location":{"latitude":45.2435799,"longitude":-123.113495},"address":{"regionCode":"us"}},{"name":"stores/7","type":"office","title":"Forest Lake Post Office","location":{"latitude":45.2722473,"longitude":-92.9854126},"address":{"regionCode":"us"}},{"name":"stores/8","type":"office","title":"Williamsfield;Williamsfield Post Office","location":{"latitude":41.5333862,"longitude":-80.5706253},"address":{"regionCode":"us"}},{"name":"stores/9","type":"office","title":"Brooklandville Post Office","location":{"latitude":39.4215546,"longitude":-76.6702347},"address":{"street":"Falls Road","city":"Brooklandville","state":"MD","zipCode":21093,"regionCode":"us"}},{"name":"stores/10","type":"office","title":"Germansville","location":{"latitude":40.7014847,"longitude":-75.7068558},"address":{"regionCode":"us"}},{"name":"stores/11","type":"office","title":"Mechanicsburg Post Office","location":{"latitude":40.2151642,"longitude":-76.9935},"address":{"street":"Simpson Ferry Road","city":"Mechanicsburg","state":"PA","zipCode":17055,"regionCode":"us"}},{"name":"stores/12","type":"office","title":"Pocahontas Trl & Heritage Mobile Village","location":{"latitude":37.2133636,"longitude":-76.6189804},"address":{"regionCode":"us"}},{"name":"stores/13","type":"office","title":"Pocahontas Trl & Carters Grove","location":{"latitude":37.2163506,"longitude":-76.621254},"address":{"regionCode":"us"}},{"name":"stores/14","type":"office","title":"Pocahontas Trl & WATA IB","location":{"latitude":37.2573624,"longitude":-76.6698685},"address":{"regionCode":"us"}},{"name":"stores/15","type":"office","title":"Pocahontas Trl & Creative Kids Daycare OB","location":{"latitude":37.2255135,"longitude":-76.626236},"address":{"regionCode":"us"}},{"name":"stores/16","type":"office","title":"MArengo Post Office","location":{"latitude":41.7998581,"longitude":-92.0712128},"address":{"regionCode":"us"}},{"name":"stores/17","type":"office","title":"Galleria Post Office","location":{"latitude":38.2518387,"longitude":-85.755928},"address":{"regionCode":"us"}},{"name":"stores/18","type":"office","title":"Fridley Post Office","location":{"latitude":45.1038361,"longitude":-93.2648544},"address":{"regionCode":"us"}},{"name":"stores/19","type":"office","title":"United States Post Office","location":{"latitude":27.9473724,"longitude":-82.4595413},"address":{"street":"North Ashley Drive","city":"Tampa","state":"FL","zipCode":33602,"regionCode":"us"}},{"name":"stores/20","type":"office","title":"Ybor City Post Office","location":{"latitude":27.9641132,"longitude":-82.4375763},"address":{"regionCode":"us"}},{"name":"stores/21","type":"office","title":"Lincoln Park Postal Store","location":{"latitude":41.9256668,"longitude":-87.6534271},"address":{"street":"North Sheffield Avenue","city":"Chicago","state":"IL","zipCode":60614,"regionCode":"us"}},{"name":"stores/22","type":"office","title":"Downtown Urbana Post Office","location":{"latitude":40.1113586,"longitude":-88.2071304},"address":{"street":"East Elm Street","city":"Urbana","state":"IL","zipCode":61801,"regionCode":"us"}},{"name":"stores/23","type":"office","title":"Vinton Post Office","location":{"latitude":37.280304,"longitude":-79.8977203},"address":{"street":"South Pollard Street","city":"Vinton","state":"VA","regionCode":"us"}},{"name":"stores/24","type":"office","title":"East Texas Post Office","location":{"latitude":40.5481148,"longitude":-75.5608597},"address":{"regionCode":"us"}},{"name":"stores/25","type":"office","title":"United States Post Office","location":{"latitude":25.9319134,"longitude":-97.4989243},"address":{"regionCode":"us"}},{"name":"stores/26","type":"office","title":"Oak View Post Office","location":{"latitude":34.3926315,"longitude":-119.300934},"address":{"regionCode":"us"}},{"name":"stores/27","type":"office","title":"Post Office","location":{"latitude":36.0843697,"longitude":-115.263725},"address":{"regionCode":"us"}},{"name":"stores/28","type":"office","title":"Main Post Office","location":{"latitude":36.0709801,"longitude":-115.139572},"address":{"regionCode":"us"}},{"name":"stores/29","type":"office","title":"Park Central Station","location":{"latitude":37.762989,"longitude":-122.242813},"address":{"regionCode":"us"}},{"name":"stores/30","type":"office","title":"US Post Office Brigantine, 08203","location":{"latitude":39.3918381,"longitude":-74.3969498},"address":{"regionCode":"us"}},{"name":"stores/31","type":"office","title":"San Juan Bautista Post Office","location":{"latitude":36.8416672,"longitude":-121.534714},"address":{"regionCode":"us"}},{"name":"stores/32","type":"office","title":"Metamora Post Office","location":{"latitude":40.7916832,"longitude":-89.3620453},"address":{"street":"East Partridge Street","city":"Metamora","state":"IL","zipCode":61548,"regionCode":"us"}},{"name":"stores/33","type":"office","title":"The Postal Store","location":{"latitude":41.5811234,"longitude":-71.5004501},"address":{"regionCode":"us"}},{"name":"stores/34","type":"office","title":"The UPS Store","location":{"latitude":42.3445663,"longitude":-71.0763855},"address":{"street":"Columbus Avenue","city":"Boston","state":"MA","zipCode":2116,"regionCode":"us"}},{"name":"stores/35","type":"office","title":"United States Post Office","location":{"latitude":39.946022,"longitude":-75.1219559},"address":{"street":"Market Street","city":"Camden","state":"NJ","zipCode":8102,"regionCode":"us"}},{"name":"stores/36","type":"office","title":"The UPS Store","location":{"latitude":34.6904373,"longitude":-82.8106232},"address":{"regionCode":"us"}},{"name":"stores/37","type":"office","title":"Canton Center Post Office","location":{"latitude":41.854126,"longitude":-72.9172897},"address":{"regionCode":"us"}},{"name":"stores/38","type":"office","title":"The UPS Store Loveland","location":{"latitude":39.2543221,"longitude":-84.2963943},"address":{"street":"Loveland Madeira Road","city":"Loveland","state":"OH","zipCode":45140,"regionCode":"us"}},{"name":"stores/39","type":"office","title":"Janesville Post Office","location":{"latitude":42.682148,"longitude":-89.0223541},"address":{"city":"Janesville","state":"WI","regionCode":"us"}},{"name":"stores/40","type":"office","title":"River Falls Post Office","location":{"latitude":44.8596268,"longitude":-92.6234512},"address":{"regionCode":"us"}},{"name":"stores/41","type":"office","title":"East Ellsworth Post Office","location":{"latitude":44.7341,"longitude":-92.4663467},"address":{"regionCode":"us"}},{"name":"stores/42","type":"office","title":"Ellsworth Post Office","location":{"latitude":44.7319107,"longitude":-92.4824142},"address":{"regionCode":"us"}},{"name":"stores/43","type":"office","title":"United States Post Office","location":{"latitude":45.9155083,"longitude":-89.2476349},"address":{"city":"Eagle River","state":"WI","zipCode":54521,"regionCode":"us"}},{"name":"stores/44","type":"office","title":"United States Post Office","location":{"latitude":45.994503,"longitude":-89.5278549},"address":{"city":"Sayner","state":"WI","zipCode":54560,"regionCode":"us"}},{"name":"stores/45","type":"office","title":"Hayward Post Office","location":{"latitude":46.0156555,"longitude":-91.4864349},"address":{"regionCode":"us"}},{"name":"stores/46","type":"office","title":"Beloit Post Office","location":{"latitude":42.4989052,"longitude":-89.0381622},"address":{"street":"Mill Street","city":"Beloit","state":"WI","zipCode":53511,"regionCode":"us"}},{"name":"stores/47","type":"office","title":"New Lisbon Post Office","location":{"latitude":43.8794136,"longitude":-90.1662369},"address":{"regionCode":"us"}},{"name":"stores/48","type":"office","title":"Plymouth Post Office","location":{"latitude":43.7488823,"longitude":-87.9770355},"address":{"regionCode":"us"}},{"name":"stores/49","type":"office","title":"Sheboygan Falls Post Office","location":{"latitude":43.7295609,"longitude":-87.8121719},"address":{"street":"Maple Street","city":"Sheboygan Falls","state":"WI","zipCode":53085,"regionCode":"us"}}],"nextPageToken":"NTA"}
```

Now try to get an individual store by changing the path to `/v1/stores/0`.

You'll get an error.
```
{"message":"UNAUTHENTICATED: Method doesn't allow unregistered callers (callers without established identity). Please use API Key or other form of API consumer identity to call this API.","code":401}
```

This is because this method requires an API key. You can get one with `gcloud`, but note that you'll need to 
change the service name from `YOUR_PROJECT` to the name of your project.

```
$ gcloud services api-keys create --key-id demo --api-target service=stores.endpoints.YOUR_PROJECT.cloud.goog 
Operation [operations/akmf.p7-622059496897-0a3a0229-45c0-484d-b770-2f63879d63cb] complete. Result: {
    "@type":"type.googleapis.com/google.api.apikeys.v2.Key",
    "createTime":"2024-10-31T15:26:42.203902Z",
    "etag":"W/\"2j5qS7GWBEUrgsfYdqi0rQ==\"",
    "keyString":"YOUR_KEY",
    "name":"projects/622059496897/locations/global/keys/demo",
    "restrictions":{
        "apiTargets":[
            {
                "service":"stores.endpoints.YOUR_PROJECT.cloud.goog"
            }
        ]
    },
    "uid":"91a05581-a418-4b80-8f5f-18fbc5ce7485",
    "updateTime":"2024-10-31T15:26:42.232775Z"
}
```

In your actual output, "YOUR_KEY" in the response above would be your actual key string.

Add it to your request in the browser by appending it to the path.
```
https://YOUR_HOST.us-west1.run.app/v1/stores/0?api_key=YOUR_KEY
```

You can also get the key to use at the command line with `gcloud` and `jq`:
```
$ KEY=$(gcloud services api-keys get-key-string projects/YOUR_PROJECT/locations/global/keys/demo --format json | jq .keyString -r)
```
