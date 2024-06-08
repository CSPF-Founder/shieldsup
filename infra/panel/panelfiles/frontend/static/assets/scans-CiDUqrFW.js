import{$ as o,r as m,s as f,a as p,b as R,l as y,c as v,h,d as D}from"./main-BCf5cIOF.js";import"./responsive.bootstrap4-4LtFXkcp.js";o(document).ready(function(){o(".table").DataTable({responsive:{details:{responsive:!0,type:"none",target:""}},order:[[0,"desc"]],language:{search:""}}),o(".dataTables_filter input").attr("placeholder","Search..."),k(),T.on("draw.dt",function(s,e){C(),u?l=b:l=g,setTimeout(S,l)});function r(s,e){var t=o(s,e.child());t.detach(),t.DataTable().destroy(),e.child.hide(),e.child.remove(),e.remove().draw(!1)}o("body").on("click",".delete-target",function(s){if(s.preventDefault(),!confirm("Are you sure want to delete the scan"))return;let e=o(this).closest("tr").attr("id");o.ajax({type:"DELETE",url:"/targets/"+e,headers:{"X-CSRF-Token":CSRF_TOKEN,"X-Requested-With":"XMLHttpRequest"},disableLoading:!0,success:function(t){if(t.redirect&&m(t.redirect),t.error)f(t.error);else if(t.success){p(t.success);let n=T.row("#"+e);r(T,n)}},error:function(t,n,a){if(t.status===422){var i=t.responseJSON;i&&i.error?f(i.error):f("Unprocessable Entity: Invalid request parameters.")}else f("An unexpected error occurred: "+n)}})})});R(function(){const r=document.getElementById("add-scan-form");if(!r)return;const s=document.getElementById("add-scan-btn");r.addEventListener("submit",function(e){s.disabled=!0,e.preventDefault(),y();const t=new FormData(r);v("/targets/add",{method:"POST",body:t}).then(n=>n.json().then(a=>({ok:n.ok,data:a}))).then(({ok:n,data:a})=>{if(h(),!n)throw new Error(a.error||"Error occurred");a.success?(D("#add-scan-form"),p(a.success)):a.redirect&&m(a.redirect),s.disabled=!1}).catch(n=>{console.log(n),h(),f(n.message),s.disabled=!1})})});const g=10*1e3,b=60*1e3;let l=b,u=!1,T,c=[];function C(){let r=!1;c=[],T.rows().every(function(){let s=this.data();const e=s.scan_status;e===TARGET_STATUS.YET_TO_START?c.push({id:s.id,scan_status:e}):e===TARGET_STATUS.SCAN_STARTED&&(c.push({id:s.id,scan_status:e}),r=!0)}),r?(l=b,u=!0):(u=!1,l=g)}function k(r){var s="/targets/list";let e={...baseAjaxData};T=o("#scan-list").DataTable({aaSorting:[],searching:!1,destroy:!0,serverSide:!0,paging:!0,initComplete:function(){o(".dt-button").removeClass("dt-button"),o(".dataTables_length").addClass("control-label pt-3")},processing:!1,responsive:{details:{responsive:!0,display:o.fn.dataTable.Responsive.display.childRowImmediate,type:"none",target:""}},serverSide:!0,stateSave:!1,lengthMenu:[[10,50,100],[10,50,100]],ajax:{url:s,dataType:"json",type:"POST",data:e,dataSrc:"records",error:function(t,n,a){t.status==303?m():f("Unable to fetch data. Please try again later. If the problem persists, please contact support.")}},columns:[{data:"target_address"},{data:"scan_status_text"},{data:"scan_started_time"},{data:"scan_completed_time"},{data:"action"}],rowId:"id"})}function x(){let r=c.map(function(s,e,t){return s.id});o.ajax({type:"POST",url:"check-multi-status",disableLoading:!0,data:{target_ids:r},success:function(s){if(s.data)for(let n in s.data){let a=s.data[n],i=a.id;var e=T.row("#"+i),t=e.data();if(t===void 0)continue;let d="";a.scan_status!==void 0&&(c=c.map(function(_,A,E){return _.id==i&&(_.scan_status=a.scan_status),_}),a.scan_status===TARGET_STATUS.SCAN_STARTED?(t.scan_started_time=a.scan_started_time,t.scan_status_text='<span class="spinner-border spinner-border-sm text-primary" aria-hidden="true"></span> <span role="status">Scanning...</span>',u=!0):t.scan_status_text=a.scan_status_text,a.scan_status==TARGET_STATUS.REPORT_GENERATED?(t.scan_completed_time=a.scan_completed_time,d='<a href="/targets/'+i+'/report" class="btn btn-sm btn-primary m-1 report-button">Report</a>',d+='<a href="/targets/'+i+'/scan-results" class="btn btn-sm btn-primary m-1 alerts-button">Alerts</a>'):(d='<a href="/targets/'+i+'/report" class="btn btn-sm btn-dark m-1 report-button disabled" disabled>Report</a>',d+='<a href="/targets/'+i+'/scan-results" class="btn btn-sm btn-dark m-1 alerts-button disabled" disabled>Alerts</a>',d='<a href="/targets/'+i+'/report" class="btn btn-sm btn-dark m-1 report-button disabled" disabled>Report</a>',d+='<a href="/targets/'+i+'/scan-results" class="btn btn-sm btn-dark m-1 alerts-button disabled" disabled>Alerts</a>'),(a.scan_status===TARGET_STATUS.SCAN_FAILED||a.scan_status===TARGET_STATUS.REPORT_GENERATED)&&(c=c.filter(function(_,A,E){return _.id!=i})),a.scan_status===TARGET_STATUS.SCAN_STARTED?d+='<a href="#" class="btn btn-sm btn-dark text-white m-1 delete-target disabled" disabled>Delete</a>':d+='<a href="#" class="btn btn-sm btn-danger text-white m-1 delete-target">Delete</a>',d!==""&&(t.action=d),e.scan_status=a.scan_status,e.data(t))}}}).always(function(){w()})}function S(){if(console.log("initiating Checking status"),c.length===0){console.log("Finished Checking status at "+new Date().toLocaleString());return}console.log("Checking status at "+new Date().toLocaleString()),x()}function w(){if(c.length===0){console.log("Finished Checking status at "+new Date().toLocaleString());return}l===g&&u?(console.log("Resetting the status check interval to 60 seconds"),l=b):(u=c.some(function(r,s,e){return r.scan_status==TARGET_STATUS.SCAN_STARTED}),!u&&l===b&&(console.log("Resetting the status check interval to 10 seconds"),l=g)),setTimeout(S,l)}
