// Part of measurement-kit <https://measurement-kit.github.io/>.
// Measurement-kit is free software. See AUTHORS and LICENSE for more
// information on the copying conditions.
#ifndef MEASUREMENT_KIT_NETTESTS_RUNNER_HPP
#define MEASUREMENT_KIT_NETTESTS_RUNNER_HPP

#include <measurement_kit/common.hpp>

namespace mk {
namespace nettests {
class NetTest;
struct RunnerCtx;

class Runner {
  public:
    Runner();
    void run_test(Var<NetTest> test, Callback<Var<NetTest>> func);
    void break_loop_();
    bool empty();
    void join_();
    ~Runner();
    static Var<Runner> global();

  private:
    Var<RunnerCtx> ctx_;
};

} // namespace nettests
} // namespace mk
#endif