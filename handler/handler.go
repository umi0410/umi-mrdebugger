package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/umi0410/umi-mrdebugger/protocol"
)

const (
	msgReadError      = "요청을 읽지 못하겠어요."
	msgParseJSONError = "요청을 JSON으로 파싱하지 못하겠어요."
	msgError          = "서비스 연결이 원활하지 않습니다."
	msgUnavailable    = "지원하지 않는 기능이에요."
	msgGreet          = "미스터 디버거 우미가 디버깅을 도와드리죠."
	msgLeave          = "남은 디버깅도 화이팅이에요!"
	msgBadIntent      = "무슨 말씀이시죠?"
)

var (
	debuggingTips = []string{
		"구글링을 해도 해도 전혀 안나오는 경우는 의외로 터무니 없는 나만의 실수일 수 있어요. 오타가 있을 지도 몰라요.",
		"서버를 로컬호스트에서 실행하신 건 아닌가요? 포트 번호도 확인해보세요!",
		"다른 환경의 서버를 찌르고 있을지도 몰라요.",
		"오타를 한 번 확인해보세요.",
		"환경변수가 아예 안 들어가고 있는지 확인해보세요.",
		"캐시 때문에 데이터가 이전 데이터를 계속 쓰고 있을 수도 있어요. 캐시를 지워보세요.",
		"브레이크 포인트를 걸고 디버그 모드로 실행해보세요. 디버그 모드는 사랑입니다.",
		"토큰이 잘 들어가고 있는지 확인해보세요.",
		"VPN을 켜야하는 작업이라면 VPN을 켰는지 확인해보세요.",
		"다른 클러스터에 작업 중이신 건 아닌가요? 사실이라면 정말 무섭네요.",
		"머리가 너무 답답할 때는 그냥 잠시 쉬고 오세요!",
	}
)

func getRandomDebuggingTips() string {
	return debuggingTips[rand.Intn(len(debuggingTips))]
}

// ServeHTTP handles CEK requests
func HelpDebugging(ctx *fiber.Ctx) error {
	var request protocol.CEKRequest

	body := ctx.Body()
	if body == nil {
		log.Error("Error during reading Request body")
		return ctx.Status(http.StatusInternalServerError).JSON(respondError(msgError))
	}

	if err := json.Unmarshal(body, &request); err != nil {
		log.Error("Error during parsing Request JSON ")
		return ctx.Status(http.StatusInternalServerError).JSON(respondError(msgError))
	}

	if protocol.CheckSignatureJinsu(ctx.GetReqHeaders()["SignatureCEK"], body) {
		log.Info("Valid request from CLOVA")
	} else {
		log.Info("Error during verifying signature")
	}

	switch request.Request.Type {
	case "LaunchRequest":
		log.Info("LaunchRequest")
		return ctx.Status(http.StatusOK).JSON(protocol.MakeCEKResponse(nil, handleLaunchRequest()))
	case "SessionEndedRequest":
		log.Info("SessionEndedRequest")
		return ctx.Status(http.StatusOK).JSON(protocol.MakeCEKResponse(nil, handleEndRequest()))
	case "IntentRequest":
		log.Info("IntentRequest")
		return ctx.Status(http.StatusOK).JSON(handleIntent(request.Request.Intent.Name))

	default:
		log.Errorf("Error wrong request type. Request: %#v", request)
		return ctx.Status(http.StatusInternalServerError).JSON(respondError(msgError))
	}
}

func Health(ctx *fiber.Ctx) error {
	return ctx.SendString("OK\n")
}

func handleLaunchRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeOutputSpeech(msgGreet + getRandomDebuggingTips()),
		ShouldEndSession: false,
	}
}

func handleEndRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeOutputSpeech(msgLeave + "\n디버깅 도움이 필요할 땐 언제든 또 미스터 디버거 우미를 찾아주세요~!"),
		ShouldEndSession: true,
	}
}

func handleIntent(intentName string) protocol.CEKResponse {
	switch intentName {
	case "AnotherTipRequested":
		return protocol.CEKResponse{
			Response: protocol.CEKResponsePayload{
				OutputSpeech:     protocol.MakeOutputSpeech(getRandomDebuggingTips()),
				ShouldEndSession: false,
			},
			SessionAttributes: nil,
		}
	case "QuitRequested":
		{
			return protocol.CEKResponse{
				Response: protocol.CEKResponsePayload{
					OutputSpeech:     protocol.MakeOutputSpeech(msgLeave),
					ShouldEndSession: true,
				},
				SessionAttributes: nil,
			}
		}
	default:
		return protocol.CEKResponse{
			Response: protocol.CEKResponsePayload{
				OutputSpeech:     protocol.MakeOutputSpeech(msgBadIntent),
				ShouldEndSession: false,
			},
			SessionAttributes: nil,
		}
	}
}

//func parseSessionAttributes(sessionAttributes map[string]string) (intent string, slots map[string]protocol.CEKSlot) {
//	slots = map[string]protocol.CEKSlot{}
//
//	for key, value := range sessionAttributes {
//		if key == "intent" {
//			intent = value
//		} else {
//			slots[key] = protocol.CEKSlot{
//				Name:  key,
//				Value: value,
//			}
//		}
//	}
//
//	return intent, slots
//}

func respondError(msg string) *protocol.CEKResponse {
	response := protocol.MakeCEKResponse(nil,
		protocol.CEKResponsePayload{
			OutputSpeech:     protocol.MakeOutputSpeech(msg),
			ShouldEndSession: true,
		})
	return &response
}

//func respondError(w http.ResponseWriter, msg string) {
//	response := protocol.MakeCEKResponse(nil,
//		protocol.CEKResponsePayload{
//			OutputSpeech:     protocol.MakeOutputSpeech(msg),
//			ShouldEndSession: true,
//		})
//
//	w.Header().Set("Content-Type", "application/json")
//	b, _ := json.Marshal(&response)
//	w.Write(b)
//}
//
//func respondSuccess(w http.ResponseWriter, response protocol.CEKResponse) {
//	w.Header().Set("Content-Type", "application/json")
//	b, _ := json.Marshal(&response)
//	w.Write(b)
//}
//
//func HealthCheck(w http.ResponseWriter, r *http.Request) {}
