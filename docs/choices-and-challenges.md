# Choices
| Date (DD/MM/YYYY) | Description of choice                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | Reason for choice                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | How was this decided?                                                                                                                         |
|-------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| 04/09/2026        | We've chosen Go as our programming language.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | Go contains everything we expect we will need in its standard library, a template engine, an easy way to set up a server, integration with many database providers.                                                                                                                                                                                                                                                                                                                                  | Majority voting in the planning channel in our Teams team.                                                                                    |
| 08/09/2026        | We have chosen to set up a separate repository for the non-legacy version of whoknows.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | In a realistic production scenario—one in which the legacy application is running, has its own CI/CD pipeline, and has actual users—it would be impossible to group the two projects into one without cumbersome logic to isolate the CI/CD pipelines and minimize the installation size on the user’s computer. Conversely, grouping them would have the advantage that rewriting would be easy, because both projects would be located close to each other and could be viewed in the same editor. | Plenary discussion (in person).                                                                                                               |
| 08/09/2026        | We have chosen the following coding conventions and will enforce them automatically in the CI/CD pipeline using `revive` and its associated GitHub Action: <ul><li>Branches are named using the format <code>[branchGroup]/branch-name</code>. Typical groups would be <code>feature</code>, <code>hotfix</code>, ...</li><li>The directory structure generally follows the Go project layout: a project in the root directory with subdirectories for specific functionality.<ul><li><code>/cmd</code> for files that contribute to executable files. Each subdirectory here should have the same name as the executable file it exposes.</li><li><code>/internal</code> for files used only internally—private libraries, in other words, files that should not be used by downstream clients.</li><li><code>/api</code> for files related to endpoints, such as the OpenAPI specification (also found in <code>/docs</code>).</li><li><code>/pkg</code> for files that can be used by downstream clients.</li><li><code>/web</code> for static assets (such as images, markup, and styling) for the site.</li><li><code>/test</code> for tests.</li></ul></li><li>Names may not be abbreviated (<code>context</code> can't become <code>ctx</code>) except for `err` for the go type `error`. We use <code>error</code> instead of <code>err</code>.</li><li>Variable names may not contain type names.</li><li>Exported variables are named in PascalCase, private ones in camelCase. The same is true of method names.</li><li>Struct names are written in PascalCase.</li><li>Indexes in loops have usual names: <code>i</code>, <code>j</code> and <code>k</code>, ...</li><li>booleans start with "is", "has" or "should", followed by a verb, noun or adjective.</li><li>"id" is written with both letters capitalized: <code>ID</code>.</li><li>Prefer code over comments, but document code in docstring equivalent comments exactly above what it documents.</li></ul>                                 | We wish to have and enforce code conventions, so that PRs violating these cannot be merged before they're corrected.                                                                                                                                                                                                                                                                                                                                                                                 | Plenary discussion following brainstorming, during which each group member proposed an idea, followed by synthesis (in person and digitally). |
| 08/09/2026        | We use the English language everywhere, except in private communication.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | Code tends to be written in English, we are to work internationally later and English allows as many people as possible to understand our code (we don't speak Spanish or Mandarin).                                                                                                                                                                                                                                                                                                                 | Plenary discussion (in person).                                                                                                               |
| 08/09/2026        | Documentation is stored in the `docs` folder at the root of our repository. This folder contains a `research` subfolder used for knowledge sharing.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      | Having a single source of truth makes the work easier, since you do not have to search for information in many different places or edit ten files every time you want to make a change. Keeping it in our repository makes it visible, which is part of Collaboration and Communication, a core value of DevOps.                                                                                                                                                                                     | Plenary discussion (in person).                                                                                                               |
| 08/09/2026        | Pull Requests (PRs) follow the following format: </p><ul><li>What does this pull request do?</li><li>What is being changed?</li><li>Why make this change/addition?</li><li>Fixes # (issue)</li><li>UI changes, if any<ul><li>Before:</li><li>After:</li></ul></li><li>Checklist before merging<ul><li><input type="checkbox" disabled> Does the code compile into a binary?</li><li><input type="checkbox" disabled> Does the code run?</li><li><input type="checkbox" disabled> Does the feature have a test?</li><li><input type="checkbox" disabled> Do all tests pass?</li><li><input type="checkbox" disabled> Does the code follow our coding conventions and use the linter?</li><li><input type="checkbox" disabled> Has the OpenAPI specification been updated if it was changed?</li></ul></li></ul> | We want a consistent format for PRs so that the review process can be formalized further by measuring the effectiveness of the procedure. We also want to be able to enforce rules that cannot be automated and remind contributors to test the code locally first, so that this becomes habitual.                                                                                                                                                                                                   | Plenary discussion following brainstorming, during which each group member proposed an idea, followed by synthesis (in person and digitally). |
| 08/09/2026        | At least one other person than the one requesting must review a PR. You may merge your own branch after a successfully completed review. A reviewer should do the following: <ol><li>Read the description, related issue, and high-level diff.<ol><li>What problem is being solved?</li><li>Is the approach appropriate (too slow/uses too much memory/bandwidth)?</li><li>Does the implementation match the requirements?</li><li>Are important edge cases or affected components missing?</li></ol></li><li>Review correctness<ol><li>Correct behavior on normal and exceptional paths</li><li>Boundary conditions and invalid input (nil, messed up structs, ...)</li><li>Error handling (is it centralized appropriately, are all potential error states covered?)</li><li>Data validation (input from forms and the like)</li><li>Database transactions (are transactions used when appropriate?)</li></ol></li><li>Review Maintainability<ol><li>Does the code conventions linter warn on their code?</li><li>Is the code understandable?</li><li>Are names precise and fitting?</li><li>Is the design consistent with the existing system?</li><li>Is complexity justified? Can anything be simplified?</li><li>Is duplication avoided without creating unnecessary abstraction?</li><li>Are comments explaining intent rather than restating the code? (if so, good)</li><li>Will another developer be able to modify this safely later?</li></ol></li><li>Review Tests<ol><li>Cover the primary behavior</li><li>Cover failure and edge cases</li><li>Are readable and deterministic</li><li>Tests behavior rather than implementation details (uses mocking in unit tests, for example)</li><li>Would fail if the feature were broken</li></ol></li></ol> | We want external verification of code to heighten quality. Peer review is a great way to get there in concert with the automated methods we will employ. | Plenary discussion (in person). | We use the library `swaggo` for semi-automatic OpenAPI spec generation from source code.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | We don't want to manually update diagrams in line with the single source of truth argument we applied to the `docs` folder.                                                                                                                                                                                                                                                                                                                                                                          | Individual research, plenary discussion (in person).                                                                                          |
| 08/09/2026        | We communicate privately in Teams, otherwise through GitHub, where we share work.  We check Teams twice per day, morning and midday, every day, except on weekends (unless we become extremely busy).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | We need to communicate asynchronously outside office hours, both to structure work and have less formal conversations.                                                                                                                                                                                                                                                                                                                                                                               | Plenary discussion (in person).                                                                                                               |
| 08/09/2026        | <ul><li>Large decisions are made through at least plenary discussion. Large decisions are decisions which direct future work like architectural decisions, code conventions, templates. Large decisions affect <em>all</em> team members. Full agreement is required.</li><li>Medium decisions may be made between two or three team members. Medium decisions relate to tasks which affect multiple, but not all people. Full agreement between those involved is required.</li><li>Small decisions may be made on one's own. Small decisions include refactoring in files, creating new ones others are not to touch at the same time.</li></ul> | We need clear thresholds for making appropriate decisions to minimize negative effects on others, and to strengthen the rationale behind large, consequential decisions. | Plenary discussion (in person).                                                                                                               |



## 17/09/2026

### Description of choice

We deployed AscendingMonk directly to our Azure VM running AlmaLinux 9.8

The current deployment is intentionally simple.

The current setup is:
* The source code is kept in our public GitHub repository.
* The repository is cloned to the VM using HTTPS.
* GitHub SSH authentication is not configured on the VM because the repository is public and therefore does not require credentials for cloning or pulling.
* Go 1.26.7, Git and GCC were installed using AlmaLinux's `dnf` package manager.
* GCC is required because the Go SQLite driver uses CGO.
* The application is compiled on the VM rather than continuously run using `go run`.
* The compiled application and its required `web/` assets are deployed under `/opt/ascendingmonk/`.
* The SQLite database is stored separately under `/var/lib/ascendingmonk/`.
* `DB_PATH` is supplied to the application as an environment variable by systemd using an absolute path.
* A dedicated Linux user called `ascendingmonk` was created to run the application.
* The application is managed by systemd and configured to start automatically when the VM starts and restart if it fails.
* The application currently listens directly on port 8080 on all VM network interfaces.
* Azure's Network Security Group allows incoming traffic to port 8080.
* `firewalld` was installed and enabled on the AlmaLinux VM as a host-level firewall.
* The public firewalld zone allows SSH and TCP port 8080. The default `dhcpv6-client` service was left enabled, while the unused Cockpit service was removed.
* SELinux remains enabled in `Enforcing` mode.
* The application currently uses HTTP - not yet HTTPS. 

The relevant deployment layout is currently:

```text
/home/<deployment-user>/AscendingMonk/
    Git repository and build environment

/opt/ascendingmonk/
    ascendingmonk
    web/
        static/
        templates/

/var/lib/ascendingmonk/
    <sqlite-database-file-with-legacy-content>.db

/etc/systemd/system/
    ascendingmonk.service
```

The deployment flow is currently approximately:

```text
Public GitHub repository
        |
        | git clone / git pull over HTTPS
        v
Source code on VM
        |
        | go build with CGO
        v
/opt/ascendingmonk/
        |
        | systemd
        v
AscendingMonk :8080
        |
        v
SQLite database in /var/lib/ascendingmonk/
```

Incoming traffic currently passes through two network filtering layers:

```text
Internet
    |
    v
Azure Network Security Group
    |
    | TCP 22 / TCP 8080
    v
AlmaLinux firewalld
    |
    | SSH / TCP 8080
    v
VM
    |
    +-- SSH :22
    |
    +-- AscendingMonk :8080
```

### Reason for choice

The deployment was deliberately kept simple because this is an early deployment exercise. 

We chose to build a binary instead of using `go run` because the Go compiler is needed during deployment, but not while the application is running. This separates building the application from running it.

A dedicated `ascendingmonk` Linux user is used rather than running the application as `<deployment-user>` or root. The service account does not need interactive login or sudo access and only needs access to the resources required by the application. This reduces the privileges available to the application.

The source code, deployed application, and mutable application data are deliberately separated:

* `/home/<deployment-user>/AscendingMonk/` is used for source code and building.
* `/opt/ascendingmonk/` contains the deployed application and its runtime web assets.
* `/var/lib/ascendingmonk/` contains persistent application data.

The SQLite database is not part of the Git repository or deployment directory.

We chose an absolute `DB_PATH` (`/var/lib/ascendingmonk/<sqlite-database-file-with-legacy-content>.db`) instead of retaining the relative path used during local development. This avoids making database access dependent on the directory from which the application happens to be started.

Systemd was chosen instead of a custom shell script, `nohup`, or keeping an SSH terminal open. It provides process supervision, logging, automatic startup after reboot, and restart behaviour using the standard service manager already provided by AlmaLinux.

The Git repository is public. Therefore, cloning over HTTPS allows the VM to pull the repository without storing a GitHub private key, personal access token, or other GitHub credentials on the VM.

Port 8080 is currently exposed directly.

We also installed `firewalld` even though Azure provides Network Security Group filtering. This gives us both cloud-level network filtering and a firewall on the VM itself, and allows us to explicitly control which services the operating system accepts from the network.

SELinux was left in `Enforcing` mode rather than disabled. The application works with SELinux enabled.

### How was this decided?

The deployment was built and tested incrementally rather than configuring everything at once.

First, we verified the environment and found that the VM was running AlmaLinux 9.8. The VM initially did not have Git, Go or GCC installed. AlmaLinux AppStream provided Go 1.26.7, which satisfies the project's `go 1.26.0` requirement, so we chose to use the distribution package rather than manually installing Go.

After installing Git, Go and GCC, we cloned the public repository over HTTPS and successfully compiled the application using CGO.

We then created the dedicated `ascendingmonk` service user and moved the SQLite database into `/var/lib/ascendingmonk/`. We tested the compiled application manually as the service user with the production `DB_PATH` before configuring systemd. This confirmed that the application, SQLite driver and database permissions worked together.

The first systemd deployment was not successful. Although the binary started, HTTP requests failed because template files such as `web/templates/layout.html` could not be found. We therefore changed the deployment layout to `/opt/ascendingmonk/`, copied both the binary and `web/` directory there, and set that directory as the systemd `WorkingDirectory`.

After this change, `curl http://localhost:8080/` returned the expected HTML.

We then verified that the application was listening on all interfaces:

```text
*:8080
```

and confirmed that both `localhost:8080` and the VM's private address `<VM private IP>:8080` returned the application.

Access through the VM's public IP initially timed out. Because the application was confirmed to work through the VM's private network interface, we identified the remaining issue as external network configuration. An Azure Network Security Group rule was added to allow inbound TCP traffic on port 8080. After this change, the application became accessible from an external browser.

Finally, `firewalld` was installed on AlmaLinux. The `public` zone was associated with `eth0`. SSH was retained, TCP port 8080 was explicitly allowed, the unused Cockpit service was removed, and the default `dhcpv6-client` service was retained. We tested both a new SSH connection and external browser access after applying the firewall configuration.

The final deployment was also tested after disconnecting the SSH sessions to verify that the application does not depend on an interactive terminal and continues to run through systemd.

## 22/09/2026
### Description of choice
We chose to implement git hooks to minimize the risk of pushing "bad" (code which the linter throws errors on or for which one or more tests fail) code to production. Two hooks were implemented, one which runs `pre-push` and another `pre-commit`. The `pre-commit` hook runs the linter on the entire codebase, and if it throws an error, git will refuse to commit the changes (unless one bypasses it manually). The one running before pushes runs all tests and the linter once more.

### Reason for choice
We opted not to run tests before each commit to maximize potential flexibility in development methodology, balancing that with our wishes for smaller commits in general. In an XP-process, a first commit for a feature may be a finished test which fails. Running tests here would be devastating in that the developer would not be allowed to commit without bypassing, which defeats the purpose of having the hook. This flexibility is crucial in this subject, since we need to have room to experiment with different methodologies to find optimal ones for our context. Our general criterion is that a feature ~ (\approx) a PR, and a feature is complete AND has tests, or has none, for tests to be written later.

### How was this decided
Personal discretion by Max-Emil. Later debate allowed for by publishing a PR, allowing for revisions and general agreement.

````markdown
## 22/09/2026

### Description of choice

We implemented a weather forecast feature as part of the Consumer Report assignment. The feature consists of a user-facing `/weather` page and an `/api/weather` endpoint. Both use the same backend weather-fetching logic.

The implementation uses Open-Meteo as the external weather provider and currently requests a seven-day forecast for a fixed location, Copenhagen. The forecast contains:

- Minimum and maximum temperature
- Precipitation probability
- Maximum wind speed
- Dominant wind direction
- WMO weather code

The HTML page converts some of the raw weather data into more readable information. For example, WMO weather codes are converted into descriptions such as "Clear sky" or "Moderate rain", and wind direction in degrees is converted into compass directions such as `N`, `SW`, and `ENE`. The JSON API retains the more machine-readable numeric values.

The integration was implemented in the backend rather than having browser-side JavaScript communicate directly with Open-Meteo.

The general flow is:

```text
Browser
   |
   | GET /weather
   v
AscendingMonk
   |
   | Check in-memory cache
   v
Cache valid? ---- yes ----> Return cached WeatherData
   |
   no
   |
   v
Open-Meteo Forecast API
   |
   | JSON
   v
Convert external response
to our WeatherData model
   |
   v
Store for 30 minutes
   |
   +----------------------+
   |                      |
   v                      v
/weather              /api/weather
HTML page             JSON response
````

A six-second HTTP timeout is used when communicating with Open-Meteo.

Successful forecasts are cached in memory for 30 minutes. The cache is protected by a mutex because multiple HTTP requests can be handled concurrently. A request that discovers an expired cache refreshes it while holding the lock, meaning other requests wait for that refresh instead of all making separate requests to Open-Meteo.

The Open-Meteo integration is tested using a fake HTTP transport instead of making real network calls during automated tests. The tests cover successful conversion of provider data, provider errors, malformed responses, inconsistent response data, request failures, timeout configuration, caching, cache expiration, and the JSON API response.

Open-Meteo attribution is displayed on the weather page.

### Reason for choice

#### Why Open-Meteo?

Open-Meteo was chosen because it was a relatively simple service to integrate for an educational, non-commercial project.

For the free/open-access service, Open-Meteo currently does not require registration or an API key for non-commercial use. Its published limits are 600 calls per minute, 5,000 per hour, 10,000 per day, and 300,000 per month. The free service has no uptime guarantee. Open-Meteo also requires attribution (a reference to where we got the data from)

These properties fit the circumstances of the project:

* We are building an educational application rather than a commercial weather product.
* We do not need to manage another secret/API key in development, CI/CD, or production.
* The API returns JSON over normal HTTP requests, making the integration relatively small.
* The API supports the daily forecast information required by the feature.
* The free usage allowance is sufficient for our expected usage, especially when combined with caching.

Avoiding an API key is useful from a DevOps perspective. An API key would have introduced additional work around:

* Secret storage
* GitHub Actions secrets
* Deployment configuration
* Key rotation
* Preventing accidental commits of credentials
* Giving developers access to development credentials

This does **not** mean that services without authentication are generally preferable. In a commercial application, a paid service with an API key, contractual terms, support, and an uptime target may be more appropriate. Open-Meteo's paid plans provide dedicated resources and a 99.9% uptime target, whereas its free API has no uptime guarantee.

Open-Meteo's current terms also distinguish between non-commercial and commercial use. A commercial deployment of AscendingMonk could therefore require us to reconsider the plan/provider rather than assuming that the current free integration can simply be carried into production.  

#### Why integrate through the backend?

We considered the integration as an architectural choice between roughly:

```text
Frontend integration

Browser -----------------> Open-Meteo
```

and:

```text
Backend integration

Browser ---> AscendingMonk ---> Open-Meteo
```

We chose the second option.

The backend approach gives AscendingMonk control over how and when the external service is contacted. This was particularly important for caching.

With a direct frontend implementation, 1,000 users loading the weather page could potentially result in approximately 1,000 requests from browsers to the weather provider:

```text
User 1  ----------------> Open-Meteo
User 2  ----------------> Open-Meteo
User 3  ----------------> Open-Meteo
...
User 1000 --------------> Open-Meteo
```

With the backend and a shared cache, many user requests can reuse one external result:

```text
User 1  ---\
User 2  ----\
User 3  -----> AscendingMonk ---> Cache ---> Open-Meteo
...         /
User 1000 -/
```

This also isolates our application from the external API's data representation. Open-Meteo returns its own JSON structure, but the rest of AscendingMonk works with our `WeatherData` and `WeatherDay` types.

Conceptually:

```text
Open-Meteo model
       |
       | conversion
       v
Our application model
       |
       +------> HTML representation
       |
       +------> JSON representation
```

This is useful because an external API should be treated as a dependency rather than as part of our own domain model. If Open-Meteo changes, or if we replace it with another provider, ideally the changes can be concentrated in the integration code instead of propagating through templates and API consumers.

A disadvantage is that our server now becomes responsible for the external dependency. A slow Open-Meteo request consumes resources on our server, and an Open-Meteo outage can affect our endpoint. This motivated the timeout and caching decisions.

#### Why use a fixed location?

The first implementation uses fixed coordinates for Copenhagen instead of allowing arbitrary locations.

This deliberately limits scope. Supporting arbitrary locations would introduce additional questions:

* How should users select a location?
* Should we accept latitude/longitude or city names?
* If city names are accepted, do we need a geocoding service?
* Should locations be validated?
* Does the cache need a separate entry for every location?
* How should cache entries be evicted?
* Could arbitrary locations dramatically increase our external API usage?

Our current cache can hold one forecast because there is only one location:

```text
weatherCache
    |
    +--> Copenhagen forecast
```

A multi-location implementation would instead require something resembling:

```text
cache
 |
 +-- Copenhagen --> forecast + expiration
 |
 +-- Aarhus -----> forecast + expiration
 |
 +-- Odense -----> forecast + expiration
 |
 +-- ...
```

At that point we would need to consider cache keys, memory growth, eviction and potentially a maximum cache size.

Open-Meteo supports requests based on latitude and longitude and can also return multiple locations in a request, so extending the feature is technically possible. ([Open-Meteo][3]) The fixed location was therefore a scope decision rather than a limitation that prevents future development.

#### Why cache weather data?

Caching was chosen primarily because the forecast does not need to be retrieved separately for every user.

Without caching:

```text
1 page request = 1 Open-Meteo request
```

Under continuous use, external API consumption would therefore increase approximately with application traffic.

With our cache:

```text
many page/API requests
        |
        v
one cached forecast
        |
refresh every 30 minutes
```

With one application instance and one fixed location, a continuously accessed application would make at most approximately:

```text
2 refreshes/hour
x 24 hours
----------------
48 refreshes/day
```

under normal successful operation.

This is far below the free service's current 10,000-call daily limit. More importantly, it demonstrates that scalability is not only about making our own server capable of receiving more requests. We also need to consider how increased traffic affects downstream dependencies.

Without caching:

```text
more users
   ↓
more requests to our server
   ↓
more requests to third party
```

With caching:

```text
more users
   ↓
more requests to our server
   ↓
mostly cache reads
   ↓
external request rate remains comparatively stable
```

The tradeoff is **freshness versus resource consumption**.

A shorter TTL gives fresher data but increases external requests:

```text
5-minute cache
= potentially 288 refreshes/day
```

A longer TTL reduces requests but can show older forecasts:

```text
6-hour cache
= potentially 4 refreshes/day
```

We selected 30 minutes as a simple compromise for this project. It is not intended as a scientifically derived optimum.

An alternative would be scheduled refreshing. Instead of refreshing when the first request arrives after expiration, a background process could fetch weather every 30 minutes:

```text
Current implementation:

request
   |
   v
cache expired?
   |
  yes
   |
   v
refresh


Alternative:

timer every 30 min
   |
   v
refresh cache

requests ---> cache
```

Scheduled refreshes can reduce latency for the first user after expiration, but they introduce background processing and continue making external requests even when nobody uses the weather feature. The request-driven TTL cache is simpler for the current application.

#### Why an in-memory cache?

The cache is stored in the Go application process rather than in SQLite or an external caching system such as Redis.

Advantages for this project include:

* Very little infrastructure.
* No additional service to deploy.
* No network call to read the cache.
* Simple implementation.
* Appropriate for one fixed location and one application instance.
* Easy to understand and test.

The main limitation is that it is **process-local**.

If the application restarts:

```text
process stops
    ↓
memory disappears
    ↓
weather cache disappears
```

The next request therefore fetches fresh weather.

More importantly, if AscendingMonk were horizontally scaled:

```text
             Load balancer
             /     |     \
            v      v      v
        Instance Instance Instance
           A        B        C
           |        |        |
        Cache A  Cache B  Cache C
```

each process would have its own cache. All three instances could independently contact Open-Meteo.

A distributed cache could instead look like:

```text
        Instance A ---\
        Instance B ----> Redis/shared cache ---> Open-Meteo
        Instance C ---/
```

That would provide shared cache state across instances, but it would also add:

* Another service to deploy
* Configuration
* Network communication
* Monitoring
* Failure modes
* Persistence/availability decisions
* Additional operational complexity

For our current single-instance educational deployment, that complexity would solve a problem we do not currently have. The in-memory cache was therefore selected as the simpler solution, while acknowledging that the decision should be revisited if the deployment architecture changes.

This illustrates an important DevOps/architecture principle: a solution should be evaluated in relation to the actual operating environment rather than automatically choosing the architecture that scales furthest.

#### Why protect the cache with a mutex?

Go HTTP handlers can execute concurrently. Therefore, multiple requests can potentially access the weather cache at the same time.

Without synchronization, operations such as:

```go
weatherCache.data = weatherData
weatherCache.expiresAt = time.Now().Add(weatherCacheDuration)
weatherCache.hasData = true
```

could occur while another goroutine reads the same state. That creates the possibility of a data race.

A mutex ensures that only one goroutine at a time enters the protected cache section:

```text
Request A ----\
Request B -----+---> mutex ---> weather cache
Request C ----/
```

Our implementation deliberately keeps the mutex locked while an expired forecast is refreshed.

That has an advantage. Imagine 100 requests arrive immediately after expiration.

Without coordination:

```text
100 requests
     |
     | cache expired
     v
100 Open-Meteo requests
```

This is sometimes called a **cache stampede** or **thundering herd**.

With the current locking strategy:

```text
Request 1 ---> obtains lock ---> Open-Meteo ---> update cache
Request 2 ---> waits
Request 3 ---> waits
...
Request 100 -> waits

after refresh:
waiting requests use cached data
```

The disadvantage is that requests waiting for the cache cannot proceed until the external request completes.

That tradeoff contributed to the decision to use an HTTP timeout.

A more advanced implementation could reduce lock contention using approaches such as:

* Read/write locks
* A dedicated refresh mechanism
* Single-flight request deduplication
* Refreshing asynchronously
* Serving stale data while a refresh occurs

Those alternatives provide more sophisticated concurrency behavior but add complexity that was not necessary for the current scope.

#### Why use an HTTP timeout?

External network requests can fail in ways that local function calls cannot.

For example:

```text
AscendingMonk ---> network ---> Open-Meteo
                      X
```

A service might be:

* Down
* Slow
* Unreachable because of routing problems
* Experiencing DNS problems
* Accepting a connection but responding very slowly

Without a timeout, our application could wait too long for the dependency.

We configured the weather HTTP client with a six-second timeout:

```go
var weatherHTTPClient = &http.Client{
    Timeout: 6 * time.Second,
}
```

This is an example of **defensive integration with an external dependency**. Our system does not control Open-Meteo, so it needs to define how long it is willing to wait.

Six seconds is a project-level engineering choice rather than a guarantee that six seconds is universally correct. A production system could derive its timeout from:

* Expected provider latency
* User-facing latency objectives
* Retry policy
* Overall request deadline
* Monitoring data
* Service-level objectives

There is also a relationship between the mutex and timeout. Because the mutex is held during a refresh, other weather requests may be waiting behind the request to Open-Meteo. The timeout therefore also provides an upper bound on how long that particular external operation is allowed to block.

#### Why separate external API structures from our application structures?

We created private structures representing Open-Meteo's response and separate public structures representing weather data in AscendingMonk.

This avoids coupling our application representation directly to the provider's JSON.

For example, conceptually:

```text
External API:

temperature_2m_max
wind_direction_10m_dominant
precipitation_probability_max

             ↓ conversion

AscendingMonk:

TemperatureMaximum
WindDirectionDominant
PrecipitationProbability
```

The provider's representation answers:

> How does Open-Meteo send this information?

Our representation answers:

> How does AscendingMonk represent weather?

Keeping those questions separate makes replacement or modification of the integration easier.

The disadvantage is additional mapping code and additional types. For a tiny integration it may look redundant, but it establishes a boundary between code we control and a contract controlled by another system.

#### Why share weather-fetching logic between HTML and API handlers?

Both:

```text
/weather
```

and:

```text
/api/weather
```

need the same forecast.

Instead of making `/weather` perform an HTTP request to our own `/api/weather` endpoint, both handlers call the same internal functionality:

```text
             fetchWeatherData()
                 /      \
                /        \
               v          v
         /weather      /api/weather
           HTML            JSON
```

An alternative would be:

```text
/weather ---> HTTP ---> /api/weather ---> Open-Meteo
```

That would introduce an unnecessary network/HTTP boundary inside the same application.

Sharing the underlying Go function means:

* The external integration exists in one place.
* Both handlers use the same cache.
* No internal HTTP serialization/deserialization is required.
* Errors are handled closer to their source.
* Tests can exercise the weather-fetching logic independently of HTTP rendering.

The API and HTML page remain different **representations** of the same underlying data.

#### Why keep numeric values in the API but format them for HTML?

The JSON API returns raw values such as the WMO weather code and wind direction in degrees, while the HTML page translates them into user-friendly descriptions.

For example:

```text
API:
weatherCode: 3
windDirectionDominant: 225

UI:
Overcast
SW
```

The purposes are different:

* An API is commonly consumed by software and benefits from structured values.
* A web page is consumed by people and benefits from readable descriptions.

Keeping presentation formatting out of the raw weather model also avoids throwing away information. A future API consumer can interpret `225` differently if needed, whereas `"SW"` has already reduced its precision.

#### Why mock Open-Meteo in automated tests?

The automated tests do not depend on real Open-Meteo responses.

Instead, the HTTP transport is replaced with a fake implementation that returns controlled responses.

Conceptually:

```text
Production:

weather code ---> HTTP client ---> Internet ---> Open-Meteo


Test:

weather code ---> HTTP client ---> fake RoundTripper
                                      |
                                      v
                              controlled JSON
```

This is particularly important from a DevOps/CI perspective.

If tests contacted the real API, a failing test could mean:

```text
our code is broken
```

but it could also mean:

```text
Internet unavailable
Open-Meteo unavailable
DNS failure
rate limit reached
provider response changed
temporary network problem
```

That would make the pipeline less deterministic.

Tests should ideally produce the same result when the application code has not changed. Removing an unnecessary external dependency from the test execution helps achieve that.

It also allows us to deliberately simulate conditions that are difficult to reproduce with the real service:

```text
HTTP 500
malformed JSON
network error
inconsistent arrays
```

For example, waiting for Open-Meteo to genuinely return malformed JSON would obviously not be a practical testing strategy.

There is nevertheless an important limitation: **mocked tests cannot prove that the real Open-Meteo integration currently works**.

We therefore used different levels of verification:

```text
Automated tests
    |
    +--> deterministic fake provider
    +--> data conversion
    +--> error handling
    +--> caching
    +--> API response

Manual/integration verification
    |
    +--> real Open-Meteo
    +--> real network
    +--> deployed environment
```

The real Open-Meteo endpoint was manually tested from the Azure-hosted environment, confirming that the deployed machine could make the required outbound HTTPS request. The completed website was also manually tested after implementation.

A more mature pipeline could add a separate integration or smoke test against the real provider. Such a test should generally be treated differently from deterministic unit tests because a third-party outage should not automatically be interpreted as a regression in our own code.

#### Why validate the external response?

Open-Meteo returns daily values in parallel arrays. Our conversion code assumes that the entries at the same index belong to the same day.

Conceptually:

```text
time:        [day1, day2, day3]
temperature: [ 15,   17,   16 ]
wind:        [  5,    7,    4 ]
```

When processing index `1`, all values should describe `day2`.

This means array lengths are an important assumption. If one array unexpectedly contained fewer values, blindly indexing it could result in incorrect processing or a runtime panic.

The integration therefore checks that the relevant arrays have consistent lengths before constructing our `WeatherDay` objects.

This is an example of treating external input as something that needs validation even when the provider is trusted. External systems are outside our control, and defensive validation can turn an unexpected response into a controlled application error.

#### Failure behavior and stale data

One tradeoff in the current implementation is what happens if the cached forecast expires and Open-Meteo cannot be reached.

The current behavior is approximately:

```text
cached forecast expires
        |
        v
try Open-Meteo
        |
        X failure
        |
        v
return error
```

An alternative would be **stale-if-error**:

```text
cached forecast expires
        |
        v
try Open-Meteo
        |
        X failure
        |
        v
serve old forecast temporarily
```

For weather data, slightly outdated information might in some circumstances be preferable to no information. However, stale data also introduces questions:

* How old may the data become?
* How do we communicate that it is stale?
* When should stale data finally stop being served?
* Should all errors allow stale data?
* How do we monitor persistent provider failures?

We kept the simpler failure behavior for this assignment. Stale-on-error would be a possible future resilience improvement rather than a requirement for the initial feature.

#### What would change in a more production-oriented deployment?

The current solution fits the circumstances of the course project, but several decisions should be reconsidered if the application requirements changed.

For example, a larger/commercial deployment might require:

* A commercial weather API licence or different provider.
* Monitoring of external API latency and errors.
* Metrics for cache hits and misses.
* Structured logging around provider failures.
* Distributed caching if multiple application instances are deployed.
* Retry logic with appropriate backoff.
* Circuit-breaking behavior during provider outages.
* Stale-on-error handling.
* Configuration rather than hard-coded location/TTL values.
* Secret management if the provider requires authentication.
* Explicit service-level objectives.
* Integration/smoke tests separate from deterministic unit tests.

This distinction is important: the choices made here are not claims that a mutex-protected in-memory cache and a hard-coded location are universally the "best" architecture. They are choices made for a small educational application, a single fixed forecast, the existing deployment setup, and the scope of the assignment.

### How was this decided?

The decision was made incrementally while implementing the Consumer Report assignment rather than by selecting the entire architecture upfront.

The process was approximately:

1. **Identify the assignment requirements.**

   * A weather forecast page was required.
   * A third-party weather service had to be investigated.
   * Cost and API-call limitations had to be considered.
   * Frontend versus backend integration had to be considered.
   * Scalability and caching had to be considered.

2. **Investigate a suitable weather provider.**

   * Open-Meteo provided the required forecast data.
   * It did not require an API key for our non-commercial educational use.
   * Its usage limits were suitable for the project.
   * Its JSON/HTTP interface could be integrated directly from Go.
   * Attribution requirements were identified and added to the page.

3. **Choose backend integration.**

   * This kept the external dependency behind AscendingMonk.
   * It allowed `/weather` and `/api/weather` to share the same integration.
   * It made server-side caching possible.
   * It avoided exposing provider-specific behavior directly to the browser.

4. **Start with a deliberately narrow data requirement.**

   * Copenhagen was selected as one fixed location.
   * A seven-day daily forecast was sufficient for the assignment.
   * Only useful daily fields were requested instead of retrieving unnecessary weather data.

5. **Create an application model separate from the provider response.**

   * Private types decode Open-Meteo JSON.
   * AscendingMonk types represent the weather data used by our own handlers.
   * HTML-specific formatting is kept in the presentation path.

6. **Expose the same data in two forms.**

   * `/api/weather` returns JSON according to the API response structure.
   * `/weather` renders the forecast for users.
   * Both call the same internal fetching logic rather than one endpoint calling the other.

7. **Test the external integration without depending on the external service.**

   * The HTTP client was made replaceable in tests.
   * A fake `RoundTripper` supplies controlled Open-Meteo responses.
   * Success and several failure scenarios were tested.
   * This makes the test suite suitable for CI because normal test execution does not require Open-Meteo to be available.

8. **Add a timeout after considering external-service failure behavior.**

   * Because an external service can become slow or unreachable, a six-second timeout was added.
   * The timeout itself was included in the automated tests.

9. **Implement the scalability plan rather than only documenting it.**

   * A 30-minute cache was added.
   * The cache reduces Open-Meteo traffic independently of the number of requests received during the cache period.
   * A mutex was added because HTTP requests can execute concurrently.
   * Tests were added to prove that a valid cache prevents another provider call and that an expired cache triggers a refresh.
   * The race detector was run to check for concurrency-related data races.

10. **Verify the complete feature outside the unit tests.**

    * The Open-Meteo request was tested from the deployed Azure environment.
    * The weather page and API were manually tested.
    * The complete automated test suite and race-detector tests passed.
    * The generated API documentation was updated.
    * The weather page was added to navigation without modifying existing elements that could interfere with the course simulation.

An important lesson from the process is that several of the non-functional requirements became visible only after the basic feature worked.

The progression was roughly:

```text
"Display weather"
       |
       v
"We need an external API"
       |
       v
"What if the API fails or is slow?"
       |
       +--> timeout
       |
       v
"What happens when many users request weather?"
       |
       +--> caching
       |
       v
"What happens when requests are concurrent?"
       |
       +--> synchronization
       |
       v
"How do we test this in CI without relying on the Internet?"
       |
       +--> mocked HTTP transport
       |
       v
"What assumptions are we making about external data?"
       |
       +--> response validation
```

This was useful from a DevOps perspective because the feature was not considered complete merely when it worked on a developer's machine. The implementation also had to consider **deployment, external dependencies, failure modes, automated testing, concurrency, API consumption, and scalability**.

The decision-making also demonstrates that architecture involves tradeoffs rather than universally correct answers. For this project we preferred relatively simple solutions:

```text
Chosen                          Alternative
----------------------------------------------------------------
Backend integration             Browser directly calls provider
In-memory cache                 Redis/shared distributed cache
30-minute request-driven TTL    Scheduled background refresh
One fixed location              Dynamic/geocoded locations
Mutex synchronization           More advanced single-flight design
Fail after refresh failure      Serve stale data on provider error
Mock provider in normal tests   Real provider in every CI run
No API key                      Authenticated commercial provider
```

The alternatives are not necessarily worse. Some would become preferable if the requirements changed. The important part of the decision is therefore not only **what we selected**, but **why that level of complexity matched the current problem** and what conditions would cause us to revisit it.


# Challenges
| Date (DD/MM/YYYY) | What happened?                                                                                                                                                                                                                                                                         | Who wrote it (and was behind it)? | What did we learn?                                                                                                                                                 |
| :---------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 02/09/2026        | I used search and replace with a naïve regex string (`print ("\w.+")`) to upgrade from Python 2 to 3. This led to repeated work, both in the code and in correcting the commit history in Git.                                                                                         | Max-Emil                          | Verify “finished” work with automated tools such as the Python interpreter, and don’t rely on `make` when it builds in release mode and suppresses error messages. |
| 09-09-2026        | Due to old habits, I looked in the list view of our GitHub Projects project and missed that Janus was already porting the about page from Python to Go. As such, we did some duplicate work. It was not all bad, since I also worked on abstractions potentially useful in later work. Janus also wrote a nice test which I hadn't done. | Max-Emil                          | Use the Kanban board for getting an overview over tasks which are in progress, reserved by others. Then decide what to work on next.                               |
