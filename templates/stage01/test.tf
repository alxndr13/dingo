resource "google_pubsub_topic" "test_topic" {
  name = "{{ .name }}"
}
