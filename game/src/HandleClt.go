package main

import (
	impl "LollipopGo/network"
	"LollipopGo2.8x/conf"
	. "LollipopGo2.8x/cxt"
	gameDB "LollipopGo2.8x/data"
	"LollipopGo2.8x/handlers/Shop"
	"LollipopGo2.8x/handlers/chat"
	. "LollipopGo2.8x/handlers/modules"
	"LollipopGo2.8x/models"
	"LollipopGo2.8x/proto/chat_proto"
	"LollipopGo2.8x/proto/comm_proto"
	"LollipopGo2.8x/proto/gm_proto"
	"LollipopGo2.8x/proto/proto_net"
	ProtoGame "LollipopGo2.8x/proto/tw_proto"
	"LollipopGo2.8x/temp"
	"fmt"
	"github.com/Golangltd/Twlib/cfgame"
	twProto "github.com/Golangltd/Twlib/proto"
	"github.com/golang/glog"
	"runtime/debug"
)

func HandleCltProtocolXL(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}) {
	// 屏蔽挂机请求打印的打印，因为太过频繁
	if v, ok := ProtocolData["Protocol2"]; ok &&
		v.(float64) != ProtoGame.C2GSUserOffLineBattleProto2 &&
		v.(float64) != proto_net.Net_HeartBeatProto {
		glog.Infof("收到数据:%+v", ProtocolData)
	}
	switch protocol {
	case float64(twProto.GGameBattleProto):
		{ // 解析战斗服 发过来的数据
			HandleCltProtocol2XL(protocol2, ProtocolData)
		}
	case float64(twProto.GGameHallProto):
		{ // 解析游戏主逻辑服务器发的数据
			HandleCltProtocol2C(protocol2, ProtocolData)
		}
	case float64(twProto.GGameConfigProto):
		{ // 获取数据库配置数据
			HandleCltProtocol2CF(protocol2, ProtocolData)
		}
	case float64(twProto.GameNetProto):
		{ // 心跳协议
			HandleCltProtocol2Net(protocol2, ProtocolData)
		}
	case float64(twProto.GGameGMProto):
		{ // 对应的GM数据：仅限测试版本
			if conf.ServerConfig().GetHotReload() {
				HandleCltProtocol2GM(protocol2, ProtocolData)
			}
		}
	default:
		glog.Info("protocol default", protocol, protocol2, ProtocolData)
	}
}

// 战斗服发过来的数据
func HandleCltProtocol2GM(protocol2 interface{}, ProtocolData map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			serer := fmt.Sprintf("%s", err)
			glog.Errorln(serer)
		}
	}()
	switch protocol2 {
	case float64(gm_proto.C2GS_GMProto2):
		{
			OpenId := ProtocolData["OpenId"].(string)
			OpType := ProtocolData["OpType"].(float64)
			//ItemId := ProtocolData["ItemId"].(float64)
			ItemNum := ProtocolData["ItemNum"].(float64)
			switch OpType {
			case 1: // 增减金币
				val, _ := M.Get(OpenId + UserKey)
				AccountId := val.(*models.Game).AccountId
				temp.GmapRoleCoin[AccountId] += int64(ItemNum)
			case 2: // 增减砖石
				val, _ := M.Get(OpenId + UserKey)
				AccountId := val.(*models.Game).AccountId
				temp.GmapRoleDiom[AccountId] += int64(ItemNum)
			case 3: // 增减道具
				{

				}
			default:
				fmt.Println("client send OpType is Wrong!")
			}
		}
	default:
		{
			glog.Info("protocol2 default")
		}
	}
}

// 战斗服发过来的数据
func HandleCltProtocol2XL(protocol2 interface{}, ProtocolData map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			serer := fmt.Sprintf("%s", err)
			glog.Errorln(serer)
		}
	}()
	switch protocol2 {
	case float64(ProtoGame.BS2GUserBattleProto2):
		BattleSendDataGameStartBattle(ConnXZ, ProtocolData)
	default:
		{
			glog.Info("protocol2 default")
		}
	}
}

func HandleCltProtocol2C(protocol2 interface{}, ProtocolData map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(fmt.Sprintf("ERROR:[%s]\nSTACK:[%s\n]", err, string(debug.Stack())))
		}
	}()

	protocolVal := protocol2.(float64)
	if HM.IsExistHandler(protocolVal) {
		HM.Handlers[protocolVal](ConnXZ, ProtocolData)
		return
	}

	// 商城系统
	if protocolVal >= ProtoGame.Proto2EquipEnd && protocolVal <= ProtoGame.Proto2EquipEnd+49 {
		Shop.HandleCltProtocol2S(ConnXZ, protocol2, ProtocolData)
		return
	}
	// 聊天系统
	if protocolVal >= chat_proto.C2GSOpenChatSysProto2 && protocolVal <= chat_proto.C2GSOpenChatSysProto2+49 {
		chat.HandleCltProtocol2Chat(ConnXZ, protocol2, ProtocolData)
		return
	}

	switch protocol2 {
	case float64(ProtoGame.C2GSUserSetSkillProto2):
		SetSkill(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserSkillUpGradeProto2):
		// SkillUpGrade(ConnXZ, ProtocolData)  服务器主动推送数据
	case float64(ProtoGame.C2GSUserActivateCardProto2):
		// 涉及表： -1- 英雄图鉴激活表  -2- 英雄卡牌图鉴表
		UserActivateGame(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserLoginProto2):
		UserLoginGame(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserChooseRoundProto2):
		UserChooseGameRound(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserStartBattleProto2):
		UserChooseStartBattle(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserChickCardProto2):
		UserClickCardGetInfo(ProtocolData)
	case float64(ProtoGame.C2GSUserRegisterProto2):
		UserRegisterInfo(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserFunctionProto2):
		UserFunctionId(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserChooseMapChapterProto2):
		UserChooseMapChapter(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserChooseChapterProto2):
	case float64(ProtoGame.C2GSUserBagProto2):
	case float64(ProtoGame.C2GSUserGetNameProto2):
	case float64(ProtoGame.C2GSUserCardUpGradeProto2):
		UserCardGrade(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserCardUpGradeBreakProto2):
		UserCardBreak(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserOffLineBattleProto2):
		UserOffLine(ProtocolData)
	case float64(ProtoGame.C2GSUserClickOnOffLineProto2):
		UserOffLineGet(ConnXZ, ProtocolData)
	case float64(comm_proto.S2GS_GetUserStProto2):
		UserStGet(ConnXZ, ProtocolData)
	case float64(comm_proto.S2GS_UpdateUserStProto2):
		UserStSet(ConnXZ, ProtocolData)
	case float64(ProtoGame.C2GSUserPlayProto2):
		UserPlayBattle(ConnXZ, ProtocolData)
	case float64(comm_proto.C2GS_GetFunctionDataProto2):
		GetRoundData(ConnXZ, ProtocolData)
	case float64(comm_proto.C2GS_UpDataItemDataProto2):
		UpdateItemData(ConnXZ, ProtocolData)
	default:
		{
			glog.Info("protocol2 default")
		}
	}
}

// 客户端发过来的数据
func HandleCltProtocol2CF(protocol2 interface{}, ProtocolData map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			serer := fmt.Sprintf("%s", err)
			glog.Errorln(serer)
		}
	}()
	switch protocol2 {
	case float64(cfgame.GetCFGameDataProto2):
		{
			tableName := ProtocolData["TableName"].(string)
			ServerId := ProtocolData["ServerId"].(string)
			data := gameDB.GetCFGameData(tableName)
			dataset := cfgame.CFRetGameData{
				Protocol:  twProto.GGameConfigProto,
				Protocol2: cfgame.CFRetGameDataProto2,
				TableName: tableName,
				Data:      data,
			}
			impl.PlayerSendMessageOfProxy(ConnXZ, dataset, ServerId)
		}
	default:
		{
			glog.Info("protocol2 default")
		}
	}
}
