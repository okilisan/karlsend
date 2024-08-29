package main

import (
	"github.com/karlsen-network/karlsend/v2/app/appmessage"
	"github.com/karlsen-network/karlsend/v2/app/protocol/common"
	"github.com/karlsen-network/karlsend/v2/domain/consensus/model/externalapi"
	"github.com/karlsen-network/karlsend/v2/domain/consensus/utils/consensushashing"
	"github.com/karlsen-network/karlsend/v2/infrastructure/network/netadapter/standalone"
	"github.com/pkg/errors"
)

func sendBlocks(address string, minimalNetAdapter *standalone.MinimalNetAdapter, blocksChan <-chan *externalapi.DomainBlock) error {
	for block := range blocksChan {
		routes, err := minimalNetAdapter.Connect(address)
		if err != nil {
			return err
		}

		blockHash := consensushashing.BlockHash(block)
		log.Infof("Sending block %s", blockHash)

		err = routes.OutgoingRoute.Enqueue(&appmessage.MsgInvRelayBlock{
			Hash: blockHash,
		})
		if err != nil {
			return err
		}

		message, err := routes.WaitForMessageOfType(appmessage.CmdRequestRelayBlocks, common.DefaultTimeout)
		if err != nil {
			return err
		}
		requestRelayBlockMessage := message.(*appmessage.MsgRequestRelayBlocks)
		if len(requestRelayBlockMessage.Hashes) != 1 || *requestRelayBlockMessage.Hashes[0] != *blockHash {
			return errors.Errorf("Expecting requested hashes to be [%s], but got %v",
				blockHash, requestRelayBlockMessage.Hashes)
		}

		err = routes.OutgoingRoute.Enqueue(appmessage.DomainBlockToMsgBlock(block))
		if err != nil {
			return err
		}

		// TODO(libp2p): Wait for reject message once it has been implemented
		err = routes.WaitForDisconnect(common.DefaultTimeout)
		if err != nil {
			return err
		}
	}
	return nil
}
