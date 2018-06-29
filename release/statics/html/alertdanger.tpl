{{define "alertdanger"}} {{if ne .alerttitle ""}}
<div class="m-section__content">
<div class="m-alert m-alert--outline alert alert-danger alert-dismissible fade show" role="alert">
    <button type="button" class="close" data-dismiss="alert" aria-label="Close"></button>
    <strong>
        {{.alerttitle}}
    </strong>
    {{.alertcontext}}
</div>
</div>
{{end}} {{end}}