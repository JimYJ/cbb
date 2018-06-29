{{define "modalBase"}}
<!--begin::Modal-->
<div class="modal fade" id="base" tabindex="-1" role="dialog" aria-labelledby="ModalTitle" aria-hidden="true">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalTitle">
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <div class="modal-body">
                <p id="modalContent">
                </p>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal">
                    取消
                </button>
                <a id="deleteit" href="">
                    <button type="button" class="btn btn-primary">
                        确定
                    </button>
                </a>

            </div>
        </div>
    </div>
</div>
<!--end::Modal-->
{{end}}