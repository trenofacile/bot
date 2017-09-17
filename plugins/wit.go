package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	messenger "github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/trenofacile/bot/witai"
)

const (
	replyUnknownCommand = "Mi spiace, non ho capito."
)

// WitAIPlugin sends the received message to witAI, gets the intent and reacts to it.
type WitAIPlugin struct {
	witAIClient *witai.Client
}

// NewWitAIPlugin returns a newly initialized WitAIPlugin with the given token.
func NewWitAIPlugin(witAIClient *witai.Client) *WitAIPlugin {
	return &WitAIPlugin{
		witAIClient: witAIClient,
	}
}

// Reply replies to a received message
func (p *WitAIPlugin) Reply(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage, mess *messenger.Messenger) {
	meaning, err := p.witAIClient.GetMeaning(msg.Text)
	if err != nil {
		log.Println(err)
		p.commandUnknown(event, opts, msg, mess)
		return
	}

	intent := p.getIntentEntityFromMeaning(meaning)

	if intent.Confidence < float64(0.7) {
		p.commandUnknown(event, opts, msg, mess)
		return
	}

	intentValue, err := intent.ValueToString()
	if err != nil {
		p.commandUnknown(event, opts, msg, mess)
		return
	}

	switch {
	case intentValue == "TrainStatus":
		p.replyWithTrainStatus(meaning, event, opts, msg, mess)
	default:
		p.commandUnknown(event, opts, msg, mess)
	}
}

// ReplyToPostback reacts to a received postback
func (p *WitAIPlugin) ReplyToPostback(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback, mess *messenger.Messenger) {
}

func (p *WitAIPlugin) commandUnknown(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage, mess *messenger.Messenger) {
	mq := messenger.MessageQuery{}
	mq.RecipientID(opts.Sender.ID)

	mq.Text(replyUnknownCommand)

	mess.SendMessage(mq)
}

type TrainStatus struct {
	StazioneUltimoRilevamento string   `json:"stazioneUltimoRilevamento,omitempty"`
	CompRitardo               []string `json:"compRitardo,omitempty"`
}

func (p *WitAIPlugin) replyWithTrainStatus(meaning *witai.Meaning, event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage, mess *messenger.Messenger) {
	client := http.Client{}
	stationID := "S09218"
	trainIDEntity, isTrainIDPresent := meaning.Entities["trainID"]
	if !isTrainIDPresent {
		p.commandUnknown(event, opts, msg, mess)
		return
	}
	trainID, err := trainIDEntity[0].ValueToString()
	if err != nil {
		p.commandUnknown(event, opts, msg, mess)
		return
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://www.viaggiatreno.it/viaggiatrenomobile/resteasy/viaggiatreno/andamentoTreno/%s/%s", stationID, trainID),
		nil,
	)
	res, err := client.Do(req)
	if err != nil {
		p.commandUnknown(event, opts, msg, mess)
		return
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		p.commandUnknown(event, opts, msg, mess)
		return
	}

	trainStatus := TrainStatus{}
	err = json.Unmarshal(data, &trainStatus)
	if err != nil {
		p.commandUnknown(event, opts, msg, mess)
		return
	}

	mq := messenger.MessageQuery{}
	mq.RecipientID(opts.Sender.ID)

	mq.Text(fmt.Sprintf("Il treno è a %s in %s", trainStatus.StazioneUltimoRilevamento, trainStatus.CompRitardo[0]))

	mess.SendMessage(mq)
}

func (p *WitAIPlugin) getIntentEntityFromMeaning(meaning *witai.Meaning) witai.Entity {
	entities, isIntentEntityArrayPresent := meaning.Entities["intent"]
	if !isIntentEntityArrayPresent {
		return witai.Entity{
			Confidence: 0,
		}
	}

	if len(entities) == 0 {
		return witai.Entity{
			Confidence: 0,
		}
	}

	return entities[0]
}
