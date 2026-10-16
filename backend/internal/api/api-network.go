package api

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linde12/gowol"

	"github.com/sarat03/lanlookout/internal/check"
	"github.com/sarat03/lanlookout/internal/portscan"
)

// getPortState godoc
// @Summary      Check port state
// @Description  Check whether a given TCP port on an address is open or closed
// @Tags         network
// @Produce      json
// @Param        addr  path      string  true  "IP address or hostname"
// @Param        port  path      string  true  "Port number"
// @Success      200   {boolean}  bool   "true if open, false if closed"
// @Router       /port/{addr}/{port} [get]
func getPortState(c *gin.Context) {
	addr := c.Param("addr")
	port := c.Param("port")
	state := portscan.IsOpen(addr, port)
	c.IndentedJSON(http.StatusOK, state)
}

// sendWOL godoc
// @Summary      Send Wake-on-LAN packet
// @Description  Send a magic packet to wake up a host by its MAC address
// @Tags         network
// @Produce      json
// @Param        mac   path      string  true  "MAC address of the host"
// @Success      200   {boolean} bool    "true if sent successfully"
// @Router       /wol/{mac} [get]
func sendWOL(c *gin.Context) {

	mac := c.Param("mac")

	packet, err := gowol.NewMagicPacket(mac)

	if !check.IfError(err) {
		err = packet.Send("255.255.255.255")

		slog.Info("Wake-on-LAN: " + mac)
	}

	c.IndentedJSON(http.StatusOK, !check.IfError(err))
}

// getIfaces godoc
// @Summary      List network interfaces
// @Description  Up, non-loopback interfaces with an IPv4 address of this machine, to pick which ones to scan
// @Tags         network
// @Produce      json
// @Success      200  {array}  string
// @Router       /ifaces [get]
func getIfaces(c *gin.Context) {
	names := []string{}
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, _ := i.Addrs()
		for _, a := range addrs { // arp-scan needs IPv4
			if n, ok := a.(*net.IPNet); ok && n.IP.To4() != nil {
				names = append(names, i.Name)
				break
			}
		}
	}
	c.IndentedJSON(http.StatusOK, names)
}
