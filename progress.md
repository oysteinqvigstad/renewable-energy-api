## Mandatory endpoints
---
The specification sheet provided for the assignment specifies four mandatory endpoints:
```
/energy/v1/renewables/current
/energy/v1/renewables/history
/energy/v1/notifications/
/energy/v1/status/
```
These have all been implemented in our software.

### Endpoint: renewables/current/ ✅
**Project specification details:** The initial endpoint focuses on returning the latest percentages of renewables in the energy mix.

**Mandatory**
- Show all countries' data by default ✅
- Allow user to specify alpha3 code ✅
- Allow user to show neighbours when specifying alpha code ✅

**Optional:** 
- Extend `{?country}` to support country name (e.g., `norway`) as input. (⚠️ *not implemented*)

### Endpoint: renewables/history/
**Project specification details:** The initial endpoint focuses on returning historical percentages of renewables in the energy mix, including individual levels, as well as mean values for individual or selections of countries.

**Mandatory**
- *Content-Type:* application/json ✅
- *Without cca:* Return mean values for all countries by default ✅
- *With cca:* Return historical data for country ✅
- Implement begin and end query options ✅

**Optional/Advanced
- Selective use of begin and end query ✅
- Time constraint for mean value when showing all countries ✅
- Sort mean values by percentage with `?sortByValue=true` ✅
- Be creative - Add whatever we feel is useful

### Endpoint: /notifications/
**Project specification details:** As an additional feature, users can register webhooks that are triggered by the service based on specified events, specifically if information about given countries (or any country) is invoked, where the minimum frequency can be specified. Users can register multiple webhooks. The registrations should survive a service restart (i.e., be persistent using a Firebase DB as backend). 

**Mandatory**
- Accept POST method on `/notifications/` for registration of webhooks ✅
- Accept DELETE method on `/notifications/{id}` for deletion of webhooks on ✅
- Accept GET method on `/notifications/{id}` for info about a single webhook ✅
- Accept GET method on `/notifications/` to list all webhooks ✅

**Advanced/optional:**
- Consider adding different notification types

### Endpoint: /status/
**Project specification details:** The status interface indicates the availability of all individual services this service depends on. These can be more services than the ones specified above (if you considered the advanced tasks). If you include more, you can specify additional keys with the suffix `api`. The reporting occurs based on status codes returned by the dependent services. The status interface further provides information about the number of registered webhooks (more details is provided in the next section), and the uptime of the service.

**Mandatory**
- Return status code 200 if everything is good ✅
- Return mandatory data ✅

## Additional requirements
---
-   All endpoints should be _tested using automated testing facilities provided by Golang_.
    -   This includes the stubbing of the third-party endpoints to ensure test reliability (removing dependency on external services).
    -   Include the testing of handlers using the httptest package. Your code should be structured to support this.
    -   Try to maximize test coverage as reported by Golang.
-   Repeated invocations for a given country and date should be cached to minimise invocation on the third-party libraries. Use Firebase for this purpose.
    -   **Advanced Task**: Implement purging of cached information for requests older than a given number of hours/days.
