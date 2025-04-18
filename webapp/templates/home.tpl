{{ template "header.tpl" . }}

{{ with .previousDate }}
    <a href="{{ addQueryParam $.CurrentURL "date" . }}" class="btn btn-primary" tabindex="-1" role="button">
        <i class="bi-arrow-left-circle-fill"></i>
    </a>
{{ end }}
{{ .item.Date }}
{{ with .nextDate }}
    <a href="{{ addQueryParam $.CurrentURL "date" . }}" class="btn btn-primary" tabindex="-1" role="button">
        <i class="bi-arrow-right-circle-fill"></i>
    </a>
{{ end }}


<main>
    {{ with .body }}
        {{ . }}
    {{ end }}
</main>

{{ template "footer.tpl" . }}