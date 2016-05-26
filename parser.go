package betterbf

import "regexp"

var emptyLine = regexp.MustCompile(` *\n+`)
var routineMatcher = regexp.MustCompile(`routine (\d)\n?`)
var tintMatcher = regexp.MustCompile(`\s*(add|sub|padd|psub|loop|end|prt|scn)\s+(.+)`)
var cronMatcher = regexp.MustCompile(`\s*_(add|sub|prt|scn|set|pset|snd|chr|if|endif|goto|exit)\s+(.+)`)
