/**
 * pingudp.go - UDP ping healthcheck
 *
 * @author Yousong Zhou <yszhou4tech@gmail.com>
 */

package healthcheck

import (
	"net"
	"time"

	"../config"
	"../core"
	"../logging"
)

const defaultUDPTimeout = 5 * time.Second

// Check executes a UDP healthcheck.
func pingudp(t core.Target, cfg config.HealthcheckConfig, result chan<- CheckResult) {
	live := false
	log := logging.For("healthcheck/pingUdp")
	defer func() {
		checkResult := CheckResult{
			Target: t,
			Live:   live,
		}
		select {
		case result <- checkResult:
		default:
			log.Warn("Channel is full. Discarding value")
		}
	}()

	addr := t.Host + ":" + t.Port
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Errorf("dial udp %s: %s", addr, err)
		return
	}
	defer conn.Close()

	timeout, _ := time.ParseDuration(cfg.Timeout)
	if timeout == time.Duration(0) {
		timeout = defaultUDPTimeout
	}
	deadline := time.Now().Add(timeout)
	err = conn.SetDeadline(deadline)
	if err != nil {
		log.Errorf("dial udp %s: set dealine: %s", addr, err)
		return
	}

	udpConn := conn.(*net.UDPConn)
	udpCfg := cfg.UdpHealthcheckConfig
	if _, err = udpConn.Write([]byte(udpCfg.Send)); err != nil {
		log.Errorf("dial udp %s: write: %s", addr, err)
		return
	}

	buf := make([]byte, len(udpCfg.Receive))
	n, _, err := udpConn.ReadFrom(buf)
	if err != nil {
		log.Errorf("dial udp %s: read: %s", addr, err)
		return
	}

	got := string(buf[0:n])
	if got != udpCfg.Receive {
		return
	}
	live = true
}
