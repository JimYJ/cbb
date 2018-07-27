{{define "modalHandelBatchVoucher"}}
<!--begin::Modal-->
<div class="modal fade" id="Batch" tabindex="-1" role="dialog" aria-labelledby="ModalAdd" aria-hidden="true">
    <div class="modal-dialog modal-lg" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAdd">
                    批量发放兑换券
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/batchvoucher">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="names" class="form-control-label">
                            兑换内容:
                        </label>
                        <input type="text" class="form-control" placeholder="兑换内容" name="content">
                    </div>
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            选择生效和失效日期:
                        </label>
                        <div class="input-daterange input-group" id="rangedate">
                            <input type="text" class="form-control m-input" name="startday" />
                            <span class="input-group-addon">
                                <i class="la la-ellipsis-h"></i>
                            </span>
                            <input type="text" class="form-control" name="endday" />
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            发放范围:
                        </label>
                        <select id="prov" onchange="showCity(this)" class="form-control m-input--fixed" name="province">
                        </select>
                        <!--城市选择-->
                        <select id="city"  class="form-control m-input--fixed" name="city"><!-- onchange="showCountry(this)" -->
                        </select>
                        <!--县区选择-->
                        <!-- <select id="country" onchange="selecCountry(this)" class="form-control m-input--fixed" name="county">
                        </select> -->
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">
                        取消
                    </button>
                    <input value="确定" type="submit" id="btns" class="btn btn-primary">
                    </input>
                </div>
            </form>
        </div>
    </div>
</div>
<!--end::Modal-->
{{end}}