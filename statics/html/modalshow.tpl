{{define "modalShow"}}
<!--begin::Modal-->
<div class="modal fade" id="base" tabindex="-1" role="dialog" aria-labelledby="ModalTitle" aria-hidden="true">
    <div class="modal-dialog modal-lg" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalTitle">
                    图片预览
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <div class="modal-body">
                <div id="showImg">
                </div>
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