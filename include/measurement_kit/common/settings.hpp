// Part of measurement-kit <https://measurement-kit.github.io/>.
// Measurement-kit is free software. See AUTHORS and LICENSE for more
// information on the copying conditions.

#ifndef MEASUREMENT_KIT_COMMON_SETTINGS_HPP
#define MEASUREMENT_KIT_COMMON_SETTINGS_HPP

#include <map>
#include <measurement_kit/common/error.hpp>
#include <sstream>
#include <string>

namespace mk {

/// This class is a specialization of the ordinary string that can also
/// store integral and boolean types. Once you have stored a value inside
/// of this class, you can access its string representation easily since
/// this class can be assigned to a string. If that fails, use the `str()`
/// method that explicitly converts to a string. To convert to integral
/// types, instead, use the `as<T>()` template method.
class SettingsEntry : public std::string {
  public:
    SettingsEntry() {}

    template <typename Type> SettingsEntry(Type value) {
        std::stringstream ss;
        ss << value;
        assign(ss.str());
    }

    template <typename Type> Type as() const {
        std::stringstream ss(c_str());
        Type value;
        ss >> value;
        if (!ss.eof()) {
            throw ValueError(); // Not all input was converted
        }
        if (ss.fail()) {
            throw ValueError(); // Input format was wrong
        }
        return value;
    }

    std::string str() const { return as<std::string>(); }

  protected:
  private:
    // NO ATTRIBUTES HERE BY DESIGN. DO NOT ADD ATTRIBUTES HERE BECAUSE
    // DOING THAT CREATES THE RISK OF OBJECT SLICING.
};

typedef std::map<std::string, SettingsEntry> Settings;

// Perhaps this could be moved in another place?
template <typename To, typename From> To lexical_cast(From f) {
    std::stringstream ss;
    To value;
    ss << f;
    ss >> value;
    if (!ss.eof()) {
        throw ValueError(); // Not all input was converted
    }
    if (ss.fail()) {
        throw ValueError(); // Input format was wrong
    }
    return value;
}

} // namespace mk
#endif
