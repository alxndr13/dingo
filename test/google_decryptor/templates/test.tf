resource "google_pubsub_topic" "test_topic" {
  name = "{{ .name }}"
}

resource "vpn_password" "secret" {
  password = "{{ .network.vpn_password }}"
}

{{ range $k, $v := .passwords }}
resource "passwords" "test_{{ $k }}" {
  value = "{{ $v }}"
}
{{- end }}
