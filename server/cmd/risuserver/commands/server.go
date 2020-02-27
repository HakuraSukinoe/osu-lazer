package commands

import (
	"github.com/deissh/osu-lazer/server/app"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Run the API server",
	RunE:         serverCmdF,
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func serverCmdF(command *cobra.Command, args []string) error {
	configPath, _ := command.Flags().GetString("config")

	server, err := app.NewServer(
		app.SetConfig(configPath),
	)
	if err != nil {
		log.Fatal().
			Err(err).
			Send()
		return err
	}
	defer server.Shutdown()

	serverErr := server.Start()
	if serverErr != nil {
		log.Fatal().
			Err(err).
			Send()
		return serverErr
	}

	notifyReady()

	interruptChan := make(chan os.Signal, 1)
	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan

	return nil
}

func notifyReady() {
	// If the environment vars provide a systemd notification socket,
	// notify systemd that the server is ready.
	systemdSocket := os.Getenv("NOTIFY_SOCKET")
	if systemdSocket != "" {
		log.Info().Msg("Sending systemd READY notification.")

		err := sendSystemdReadyNotification(systemdSocket)
		if err != nil {
			log.Fatal().
				Err(err).
				Send()
		}
	}
}

func sendSystemdReadyNotification(socketPath string) error {
	msg := "READY=1"
	addr := &net.UnixAddr{
		Name: socketPath,
		Net:  "unixgram",
	}
	conn, err := net.DialUnix(addr.Net, nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(msg))
	return err
}
