<!doctype html>
<html lang="{{.languageCode}}">
<head>
<meta charset="utf-8">
<title>{{if (eq .newUrlPath "index.html")}}{{.homeTitle}}{{else}}{{.dataTitle}}{{.suffixTitle}}{{end}}</title>
<meta name="viewport" content="width=device-width,initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no">
<meta name="format-detection" content="telephone=no,email=no,address=no">
<link rel="stylesheet" href="https://cdn.bootcss.com/semantic-ui/2.2.13/semantic.min.css">
<link rel="stylesheet" href="{{.fixLink}}asset/css/style.css">
</head>
<body>

<div class="ui left vertical menu sidebar">
    <div class="menu">
        {{.dataMenu}}

        <div></div>
        <ul class="made-by">
            <li><a href="https://easydoc.089858.com" target="_blank" title="EasyDoc">EasyDoc</a></li>
        </ul>
    </div>
</div>

<div class="pusher">
    <div class="ui vertical">
        <div class="ui fixed inverted borderless menu">
            <a href="javascript:;" class="item" id="btn-sidebar"><i class="sidebar icon"></i></a>
            <a href="{{.fixLink}}index.html" class="item">Home</a>
            <div class="right menu">
                <a  href="https://easydoc.089858.com" class="item" target="_blank" title="EasyDoc">EasyDoc</a>
            </div>
        </div>

        <div class="ui grid new-grid">
            <div class="sixteen wide column">
                <div class="ui raised segment">
                    <strong class="ui teal ribbon label">{{.dataTitle}}</strong>
                    <div class="content">
                        {{.dataDoc}}
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<p id="back2top">&and;</p>

<script src="https://cdn.bootcss.com/jquery/2.2.3/jquery.min.js"></script>
<script src="https://cdn.bootcss.com/semantic-ui/2.2.13/semantic.min.js"></script>
<script src="https://cdn.bootcss.com/jquery-throttle-debounce/1.1/jquery.ba-throttle-debounce.min.js"></script>
<script src="{{.fixLink}}asset/js/app.js"></script>
</body>
</html>
