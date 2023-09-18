provider "kubernetes" {
  config_path = "~/.kube/config"
}

data "kubernetes_service" "nginx_service" {
  metadata {
    name      = "ingress-nginx-controller"
    namespace = helm_release.nginx_ingress.namespace
  }


}

output "load_balancer_address" {
  description = "Load balancer URL"
  value       = data.kubernetes_service.nginx_service.status.0.load_balancer.0.ingress.0.hostname
}

data "aws_route53_zone" "sales-dev" {
  name         = "sales-dev.codefresh.io"
}

output "zone_id" {
  description = "Zone ID"
  value       = data.aws_route53_zone.sales-dev.zone_id
}

resource "aws_route53_record" "multi-demo" {
  zone_id = data.aws_route53_zone.sales-dev.zone_id
  name    = "multi.sales-dev.codefresh.io"
  type    = "CNAME"
  ttl     = "300"
  records = [data.kubernetes_service.nginx_service.status.0.load_balancer.0.ingress.0.hostname]
}

resource "aws_route53_record" "star-multi-demo" {
  zone_id = data.aws_route53_zone.sales-dev.zone_id
  name    = "*.multi.sales-dev.codefresh.io"
  type    = "CNAME"
  ttl     = "300"
  records = [data.kubernetes_service.nginx_service.status.0.load_balancer.0.ingress.0.hostname]
}