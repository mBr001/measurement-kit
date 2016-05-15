// Part of measurement-kit <https://measurement-kit.github.io/>.
// Measurement-kit is free software. See AUTHORS and LICENSE for more
// information on the copying conditions.
#ifndef MEASUREMENT_KIT_MLABNS_MLABNS_HPP
#define MEASUREMENT_KIT_MLABNS_MLABNS_HPP

#include <cstddef>
#include <functional>
#include <initializer_list>
#include <measurement_kit/common.hpp>
#include <string>
#include <utility>
#include <vector>

namespace mk {
namespace mlabns {

MK_DEFINE_ERR(5000, InvalidPolicyError, "unknown_failure 5000")
MK_DEFINE_ERR(5001, InvalidAddressFamilyError, "unknown_failure 5001")
MK_DEFINE_ERR(5002, InvalidMetroError, "unknown_failure 5002")
MK_DEFINE_ERR(5003, InvalidToolNameError, "unknown_failure 5003")
MK_DEFINE_ERR(5004, UnexpectedHttpStatusCodeError, "unknown_failure 5004")

/// Reply to mlab-ns query.
class Reply {
  public:
    std::string city;            ///< City where sliver is.
    std::string url;             ///< URL to access sliver using HTTP.
    std::vector<std::string> ip; ///< List of IP addresses of sliver.
    std::string fqdn;            ///< FQDN of sliver.
    std::string site;            ///< Site where sliver is.
    std::string country;         ///< Country where sliver is.
};

/// Query mlab-ns and receive response.
void query(std::string tool, Callback<Error, Reply> callback,
           Settings settings = {}, Var<Reactor> reactor = Reactor::global(),
           Var<Logger> logger = Logger::global());

} // namespace mlabns
} // namespace mk
#endif
