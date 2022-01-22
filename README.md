# Bhojpur Pay - Data Processing Engine

The Bhojpur Pay is a software-as-a-service product used as a Payments Data Processing Engine based on Bhojpur.NET Platform for application delivery.

## 3rd party Gateway Integration

The Bhojpur Pay also features integration with some third-party payment gateways. There is an abstracted payment interface for software developers. It provides a unified API for different payment gateways.

### Usage

```go
import "github.com/bhojpur/pay/pkg/gateway/stripe"

func main() {
  Stripe := stripe.New(&stripe.Config{
    Key: config.Key,
  })
}
```
