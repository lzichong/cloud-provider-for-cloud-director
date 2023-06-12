/*
 * VMware Cloud Director OpenAPI
 *
 * VMware Cloud Director OpenAPI is a new API that is defined using the OpenAPI standards.<br/> This ReSTful API borrows some elements of the legacy VMware Cloud Director API and establishes new patterns for use as described below. <h4>Authentication</h4> Authentication and Authorization schemes are the same as those for the legacy APIs. You can authenticate using the JWT token via the <code>Authorization</code> header or specifying a session using <code>x-vcloud-authorization</code> (The latter form is deprecated). <h4>Operation Patterns</h4> This API follows the following general guidelines to establish a consistent CRUD pattern: <table> <tr>   <th>Operation</th><th>Description</th><th>Response Code</th><th>Response Content</th> </tr><tr>   <td>GET /items<td>Returns a paginated list of items<td>200<td>Response will include Navigational links to the items in the list. </tr><tr>   <td>POST /items<td>Returns newly created item<td>201<td>Content-Location header links to the newly created item </tr><tr>   <td>GET /items/urn<td>Returns an individual item<td>200<td>A single item using same data type as that included in list above </tr><tr>   <td>PUT /items/urn<td>Updates an individual item<td>200<td>Updated view of the item is returned </tr><tr>   <td>DELETE /items/urn<td>Deletes the item<td>204<td>No content is returned. </tr> </table> <h5>Asynchronous operations</h5> Asynchronous operations are determined by the server. In those cases, instead of responding as described above, the server responds with an HTTP Response code 202 and an empty body. The tracking task (which is the same task as all legacy API operations use) is linked via the URI provided in the <code>Location</code> header.<br/> All API calls can choose to service a request asynchronously or synchronously as determined by the server upon interpreting the request. Operations that choose to exhibit this dual behavior will have both options documented by specifying both response code(s) below. The caller must be prepared to handle responses to such API calls by inspecting the HTTP Response code. <h5>Error Conditions</h5> <b>All</b> operations report errors using the following error reporting rules: <ul>   <li>400: Bad Request - In event of bad request due to incorrect data or other user error</li>   <li>401: Bad Request - If user is unauthenticated or their session has expired</li>   <li>403: Forbidden - If the user is not authorized or the entity does not exist</li> </ul> <h4>OpenAPI Design Concepts and Principles</h4> <ul>   <li>IDs are full Uniform Resource Names (URNs).</li>   <li>OpenAPI's <code>Content-Type</code> is always <code>application/json</code></li>   <li>REST links are in the Link header.</li>   <ul>     <li>Multiple relationships for any link are represented by multiple values in a space-separated list.</li>     <li>Links have a custom VMware Cloud Director-specific &quot;model&quot; attribute that hints at the applicable data         type for the links.</li>     <li>title + rel + model attributes evaluates to a unique link.</li>     <li>Links follow Hypermedia as the Engine of Application State (HATEOAS) principles. Links are present if         certain operations are present and permitted for the user&quot;s current role and the state of the         referred entities.</li>   </ul>   <li>APIs follow a flat structure relying on cross-referencing other entities instead of the navigational style       used by the legacy VMware Cloud Director APIs.</li>   <li>Most endpoints that return a list support filtering and sorting similar to the query service in the legacy       VMware Cloud Director APIs.</li>   <li>Accept header must be included to specify the API version for the request similar to calls to existing legacy       VMware Cloud Director APIs.</li>   <li>Each feature has a version in the path element present in its URL.<br/>       <b>Note</b> API URL's without a version in their paths must be considered experimental.</li> </ul> 
 *
 * API version: 37.2
 * Contact: https://code.vmware.com/support
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

// Describes a Firewall rule for an edge gateway backed by NSX-T.
type EdgeFirewallRule struct {
	// The unique id of this firewall rule. If a rule with the ruleId is not already present, a new rule will be created. If it already exists, the rule will be updated. 
	Id string `json:"id,omitempty"`
	// Name for the rule.
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	// List of source groups for firewall rule. It specifies the sources of network traffic for the firewall rule. Null value or an empty list will be treated as \"ANY\" which means traffic from any source. For Distributed Firewall rules, an entry with an id of urn:vcloud:firewallGroup:internal can be used to specify all internal vDC Group network traffic. 
	SourceFirewallGroups []EntityReference `json:"sourceFirewallGroups,omitempty"`
	// List of source groups for firewall rule. It specifies the destinations of network traffic for the firewall rule. Null value or an empty list will be treated as \"ANY\" which means traffic to any destination. For Distributed Firewall rules, an entry with an id of urn:vcloud:firewallGroup:internal can be used to specify all internal vDC Group network traffic. 
	DestinationFirewallGroups []EntityReference `json:"destinationFirewallGroups,omitempty"`
	// The list of application ports where this firewall rule is applicable. Null value or an empty list will be treated as \"ANY\" which means rule applies to all ports. 
	ApplicationPortProfiles []EntityReference `json:"applicationPortProfiles,omitempty"`
	// Type of IP packet that should be matched while enforcing the rule. Default value is IPV4_IPV6. 
	IpProtocol *FirewallRuleIpProtocol `json:"ipProtocol,omitempty"`
	// The action to be applied to all the traffic that meets the firewall rule criteria. It determines if the rule permits or blocks traffic. This property is now deprecated and replaced with actionValue. Property is required if actionValue is not set. 
	Action *FirewallRuleAction `json:"action,omitempty"`
	// The action to be applied to all the traffic that meets the firewall rule criteria. It determines if the rule permits or blocks traffic. Property is required if action is not set. Below are valid values. <ul>   <li> <code> ALLOW </code> permits traffic to go through the firewall.   <li> <code> DROP </code> blocks the traffic at the firewall. No response is sent back to the source.   <li> <code> REJECT </code> blocks the traffic at the firewall. A response is sent back to the source. </ul> 
	ActionValue string `json:"actionValue,omitempty"`
	// Specifies the direction of the network traffic. Default value is IN_OUT. 
	Direction *FirewallRuleDirection `json:"direction,omitempty"`
	// Whether packet logging is enabled for firewall rule.
	Logging bool `json:"logging,omitempty"`
	// The list of layer 7 network context profiles where this firewall rule is applicable. Null value or an empty list will be treated as \"ANY\" which means rule applies to all applications and domains. 
	NetworkContextProfiles []EntityReference `json:"networkContextProfiles,omitempty"`
	// Whether the firewall rule is enabled.
	Enabled bool           `json:"enabled,omitempty"`
	Version *ObjectVersion `json:"version,omitempty"`
	// Text for user entered comments on the firewall rule. Length cannot exceed 2048 characters.
	Comments string `json:"comments,omitempty"`
	// Used to limit application of this firewall rule to the specified Org VDC or segment backed external network. Only networks connected as service interfaces are usable.
	AppliedTo *EntityReference `json:"appliedTo,omitempty"`
}