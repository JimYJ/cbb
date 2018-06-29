{{define "modalHandelGoodsRedeem"}}
<!--begin::Modal-->
<div class="modal fade" id="AddEdit" tabindex="-1" role="dialog" aria-labelledby="ModalAdd" aria-hidden="true">
    <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAdd">
                    新增记录
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/addvendor">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="names" class="form-control-label">
                            蝴蝶类型:
                        </label>
                        <select name="butterflyid" class="form-control m-input--fixed" id="butterflyid">
                            {{range $b:=.butterflylist}}
                            <option value="{{$b.id}}">
                                {{$b.name}}
                            </option>
                            {{end}}
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="numbers" class="form-control-label">
                            数量:
                        </label>
                        <input type="text" class="form-control" placeholder="数量" id="numbers" name="numbers">
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">
                        取消
                    </button>
                    <input value="确定" type="submit" class="btn btn-primary">
                    </input >
                </div>
            </form>
        </div>
    </div>
</div>
<!--end::Modal-->
{{end}}