{{ template "header.tpl" . }}

<div class="container-fluid">
    <div class="row justify-content-center">
        <div class="col-lg-8 col-md-10 col-sm-12">
            <!-- Date Navigation -->
            <div class="d-flex justify-content-between align-items-center mb-4 mt-3">
                {{ with .previousDate }}
                    <a href="{{ addQueryParam $.CurrentURL "date" . }}" class="btn btn-outline-primary" title="Previous Entry">
                        <i class="bi bi-arrow-left-circle"></i> Previous
                    </a>
                {{ else }}
                    <div></div>
                {{ end }}
                
                <div class="text-center">
                    <span class="badge bg-secondary fs-6 px-3 py-2">
                        <i class="bi bi-calendar3"></i> {{ .item.Date }}
                    </span>
                </div>
                
                {{ with .nextDate }}
                    <a href="{{ addQueryParam $.CurrentURL "date" . }}" class="btn btn-outline-primary" title="Next Entry">
                        Next <i class="bi bi-arrow-right-circle"></i>
                    </a>
                {{ else }}
                    <div></div>
                {{ end }}
            </div>

            <!-- Main Diary Entry Card -->
            <div class="card shadow-sm mb-4">
                {{ if .item.Title }}
                    <div class="card-header bg-primary text-white">
                        <h1 class="card-title mb-0 h3">
                            <i class="bi bi-journal-text me-2"></i>{{ .item.Title }}
                        </h1>
                    </div>
                {{ end }}
                
                <div class="card-body">
                    {{ if not .item.Title }}
                        <div class="text-center text-muted mb-3">
                            <i class="bi bi-journal-text fs-1"></i>
                            <p class="mt-2">Daily Entry</p>
                        </div>
                    {{ end }}
                    
                    <main class="diary-content">
                        {{ with .body }}
                            {{ . }}
                        {{ else }}
                            <div class="text-center text-muted py-5">
                                <i class="bi bi-pencil-square fs-1 mb-3"></i>
                                <p class="fs-5">No content yet for this day</p>
                                <a href="/web/edit?date={{ .item.Date }}" class="btn btn-primary">
                                    <i class="bi bi-plus-circle me-1"></i>Add Entry
                                </a>
                            </div>
                        {{ end }}
                    </main>
                </div>
                
                {{ if .item.Tags }}
                    <div class="card-footer bg-light">
                        <div class="d-flex flex-wrap align-items-center">
                            <span class="text-muted me-2">
                                <i class="bi bi-tags"></i> Tags:
                            </span>
                            {{ range .item.Tags }}
                                {{ if . }}
                                    <span class="badge bg-info text-dark me-1 mb-1">{{ . }}</span>
                                {{ end }}
                            {{ end }}
                        </div>
                    </div>
                {{ end }}
            </div>

            <!-- Quick Actions -->
            <div class="d-flex justify-content-center gap-2 mb-4">
                <a href="/web/edit?date={{ .item.Date }}" class="btn btn-outline-primary">
                    <i class="bi bi-pencil-square me-1"></i>Edit Entry
                </a>
                <a href="/web/edit" class="btn btn-outline-success">
                    <i class="bi bi-plus-circle me-1"></i>New Entry
                </a>
            </div>
        </div>
    </div>
</div>

<style>
    .diary-content {
        line-height: 1.6;
        font-size: 1.1rem;
    }
    
    .diary-content h1, .diary-content h2, .diary-content h3 {
        color: #495057;
        margin-top: 1.5rem;
        margin-bottom: 1rem;
    }
    
    .diary-content p {
        margin-bottom: 1rem;
    }
    
    .diary-content ul, .diary-content ol {
        margin-left: 1.5rem;
        margin-bottom: 1rem;
    }
    
    .diary-content blockquote {
        border-left: 4px solid #007bff;
        padding-left: 1rem;
        margin: 1rem 0;
        font-style: italic;
        color: #6c757d;
    }
    
    .card {
        border: none;
        border-radius: 12px;
    }
    
    .card-header {
        border-radius: 12px 12px 0 0 !important;
    }
    
    .badge {
        border-radius: 8px;
    }
</style>

{{ template "footer.tpl" . }}