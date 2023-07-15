package lib

import (
	"fmt"
	"io"

	"github.com/samyfodil/taubyte-llama-satellite/sdk"
)

type question struct {
	ctx    string
	prompt string
}

var questions = []*question{
	{
		ctx: `Taubyte is a set of tools designed to empower you as a developer. It aims to accelerate the adoption of Web 3 and Edge Computing technologies, lower their associated costs, and promote innovation by creating a universe of interoperable cloud constellations. With Taubyte, you can navigate the future of technology with confidence and ease.
Acting as a cloud-native layer for Web 3 and Edge Computing, it brings the familiar, developer-friendly model of cloud computing to these rapidly evolving areas of technology.
As a developer, you often face the challenges of learning and adapting to new paradigms, tools, and frameworks that can be complex and time-consuming. Taubyte aims to simplify this process by providing a straightforward and intuitive interface that bridges the gap between traditional cloud computing and the world of Web 3 and Edge Computing.
But Taubyte is not just about making life easier for developers. One of the significant hurdles in the industry today is the underutilization of infrastructure. By optimizing resource usage, Taubyte helps to reduce costs and enhance efficiency, making it a cost-effective solution for application development and deployment.
Furthermore, Taubyte promotes interoperability across various infrastructures. This characteristic is crucial in today’s tech ecosystem, where being able to integrate and communicate smoothly between different platforms can significantly enhance the value delivered to end-users and clients.`,
		prompt: "What is Taubyte?",
	}, {
		prompt: "what is 1 * 2 + 325 ?",
	}, {
		prompt: "what is 333 + 25 ?",
	}, {
		prompt: "what is 333 / 3 ?",
	}, {
		prompt: "Give me five prime numbers",
	},
}

func (q *question) String() string {
	return q.ctx + "\n" + q.prompt
}

//export wapredict
func wapredict(qi uint32) uint32 {
	question := questions[qi]
	fmt.Println(question.ctx)
	fmt.Println("Q: " + question.prompt)
	fmt.Print("A: ")

	p, err := sdk.Predict(
		question.String(),
		sdk.WithTopK(90),
		sdk.WithTopP(0.86),
		sdk.WithBatch(5),
		sdk.WithPenalty(1),
	)
	if err != nil {
		panic(err)
	}

	for {
		token, err := p.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Print(token)
	}

	fmt.Println("\n--")

	return 0
}
