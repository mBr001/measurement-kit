// Part of measurement-kit <https://measurement-kit.github.io/>.
// Measurement-kit is free software. See AUTHORS and LICENSE for more
// information on the copying conditions.

#ifndef SRC_OONI_TCP_TEST_HPP
#define SRC_OONI_TCP_TEST_HPP

#include <measurement_kit/net/buffer.hpp>

#include <measurement_kit/common/logger.hpp>
#include <measurement_kit/common/settings.hpp>
#include <measurement_kit/common/var.hpp>

#include "src/ooni/ooni_test_impl.hpp"
#include "src/net/connection.hpp"

namespace mk {

namespace net {
class Connection;
}

namespace ooni {

typedef Var<net::Connection> TCPClient; /* XXX */

class TCPTestImpl : public ooni::OoniTestImpl {
    using ooni::OoniTestImpl::OoniTestImpl;

  public:
    TCPTestImpl(std::string input_filepath_, Settings options_)
        : ooni::OoniTestImpl(input_filepath_, options_) {
        test_name = "tcp_test";
        test_version = "0.0.1";
    };

    void connect(Settings options, std::function<void(TCPClient)> &&cb) {

        if (options["port"] == "") {
            throw std::runtime_error("Port is required");
        }
        if (options["host"] == "") {
            options["host"] = "localhost";
        }

        auto connection = std::make_shared<net::Connection>(
            "PF_UNSPEC", options["host"].c_str(), options["port"].c_str(),
            &logger, poller);

        //
        // FIXME The connection and this are bound in the
        // callbacks below, but they have possibly different
        // life cycles, which is &disaster.
        //

        connection->on_error([cb, this, connection](Error e) {
            connection->on_error(nullptr);
            connection->on_connect(nullptr);
            entry["error_code"] = (int)e;
            entry["connection"] = "failed";
            cb(connection);
        });
        connection->on_connect([this, cb, connection]() {
            connection->on_error(nullptr);
            connection->on_connect(nullptr);
            entry["connection"] = "success";
            cb(connection);
        });
    }
};

} // namespace ooni
} // namespace mk
#endif
