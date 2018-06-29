{{define "leftmenu"}}
<!-- BEGIN: Left Aside -->
<button class="m-aside-left-close  m-aside-left-close--skin-dark " id="m_aside_left_close_btn">
    <i class="la la-close"></i>
</button>
<div id="m_aside_left" class="m-grid__item	m-aside-left  m-aside-left--skin-dark ">
    <!-- BEGIN: Aside Menu -->
<div 
id="m_ver_menu" 
class="m-aside-menu  m-aside-menu--skin-dark m-aside-menu--submenu-skin-dark " 
data-menu-vertical="true"
data-menu-scrollable="false" data-menu-dropdown-timeout="500"  
>
        <ul class="m-menu__nav  m-menu__nav--dropdown-submenu-arrow ">
            <li class="m-menu__item  m-menu__item--active" aria-haspopup="true" >
                <a  href="/" class="m-menu__link ">
                    <i class="m-menu__link-icon flaticon-line-graph"></i>
                    <span class="m-menu__link-title">
                        <span class="m-menu__link-wrap">
                            <span class="m-menu__link-text">
                                管理后台
                            </span>
                            <!-- <span class="m-menu__link-badge">
                                <span class="m-badge m-badge--danger">
                                    2
                                </span>
                            </span> -->
                        </span>
                    </span>
                </a>
            </li>
            <li class="m-menu__section">
                <h4 class="m-menu__section-text">
                    菜单
                </h4>
                <i class="m-menu__section-icon flaticon-more-v3"></i>
            </li>
            {{range $v:= .menu}}
            <li class="m-menu__item  m-menu__item--submenu" id="masterMenu" onclick="getShowMenu('{{$v.id}}')" aria-haspopup="true" mid="{{$v.id}}"  data-menu-submenu-toggle="hover">
                <a href="{{$v.path}}" class="m-menu__link m-menu__toggle" >
                    <i class="m-menu__link-icon {{$v.icon}}"></i>
                    <span class="m-menu__link-text">
                        {{$v.name}}
                    </span>
                    <i class="m-menu__ver-arrow la la-angle-right"></i>
                </a>
                <div class="m-menu__submenu">
                    <span class="m-menu__arrow"></span>
                    <ul class="m-menu__subnav">
                            {{range $v2 := $v.list}}
                        <li class="m-menu__item " aria-haspopup="true" >
                                <a  href="{{$v2.path}}" class="m-menu__link ">
                                    <i class="m-menu__link-bullet m-menu__link-bullet--dot">
                                        <span></span>
                                    </i>
                                    <span class="m-menu__link-text">
                                            {{$v2.name}}
                                    </span>
                                </a>
                            </li>
                        {{end}}
                    </ul>
                </div>
            </li>
            {{end}}
        </ul>
    </div>
    <!-- END: Aside Menu -->
</div>
<!-- END: Left Aside -->
{{end}}