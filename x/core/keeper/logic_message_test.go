package keeper_test

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_message.go

* DispatchMessage (invalid) with non-existing Mailbox ID
* ProcessMessage (invalid) with non-existing Mailbox ID
* ProcessMessage (invalid) with invalid hex message
* ProcessMessage (invalid) already processed message (replay protection)
* ProcessMessage (invalid) with invalid message: non-registered recipient
* ProcessMessage (invalid) with invalid message: invalid version

*/

var _ = Describe("logic_message.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress
	var sender i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		sender = i.GenerateTestValidatorAddress("Sender")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("DispatchMessage (invalid) with non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		hexSender, _ := util.DecodeHexAddress(sender.Address)
		recipient, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		body, _ := hex.DecodeString("0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b")

		_, err = s.App().HyperlaneKeeper.DispatchMessage(
			s.Ctx(),
			nonExistingMailboxId,
			hexSender,
			sdk.NewCoins(sdk.NewCoin("acoin", math.NewInt(1000000))),
			1,
			recipient,
			body,
			util.StandardHookMetadata{
				GasLimit: math.NewInt(50000),
				Address:  sender.AccAddress,
			},
			nil,
		)

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))

		verifyDispatch(s, mailboxId, 0)
	})

	It("ProcessMessage (invalid) with non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      1337,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: nonExistingMailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	It("ProcessMessage (invalid) with wrong destination domain", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      1337,
			Sender:      senderHex,
			Destination: 2,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("message destination %v does not match local domain %v", 2, 1)))
	})

	It("ProcessMessage (invalid) with invalid hex message", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String()[:util.BodyOffset-1],
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid hyperlane message"))
	})

	It("ProcessMessage (invalid) already processed message (replay protection)", func() {
		// Arrange
		mailboxId, _, _, ismId := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		// Register a mock recipient
		mockApp := i.CreateMockApp(s.App().HyperlaneKeeper.AppRouter())
		recipient, err := mockApp.RegisterApp(s.Ctx(), ismId)
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipient,
			Body:        []byte(""),
		}

		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})
		Expect(err).To(BeNil())

		// Expect our mock app to have been called
		callcount, message, mailboxId := mockApp.CallInfo()
		Expect(callcount).To(Equal(1))
		Expect(message.String()).To(Equal(message.String()))
		Expect(mailboxId.String()).To(Equal(mailboxId.String()))
		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("already received messsage with id %s", hypMsg.Id())))

		// Expect our mock app to not have been called again
		callcount, _, _ = mockApp.CallInfo()
		Expect(callcount).To(Equal(1))
	})

	It("ProcessMessage (invalid) with invalid message: non-registered recipient", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("id %v not found", recipientHex)))
	})

	It("ProcessMessage (invalid) with invalid message: invalid version", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		var version uint8 = 2

		hypMsg := util.HyperlaneMessage{
			Version:     version,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("unsupported message version %v", version)))
	})
})
