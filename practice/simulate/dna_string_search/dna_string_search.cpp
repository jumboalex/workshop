#include <iostream>
#include <string>
#include <vector>
#include <unordered_map>
#include <unordered_set>

using namespace std;

// Global map for IUPAC nucleotide codes
// Maps each pattern character to the set of bases it can match
unordered_map<char, unordered_set<char>> patternCharMap;

// Initialize the IUPAC nucleotide mapping
void initializePatternMap() {
    // Standard bases
    patternCharMap['A'] = {'A'};
    patternCharMap['C'] = {'C'};
    patternCharMap['G'] = {'G'};
    patternCharMap['T'] = {'T'};

    // Degenerate bases
    patternCharMap['R'] = {'A', 'G'};  // puRine
    patternCharMap['Y'] = {'C', 'T'};  // pYrimidine
    patternCharMap['M'] = {'A', 'C'};  // aMino
    patternCharMap['K'] = {'G', 'T'};  // Keto
    patternCharMap['W'] = {'A', 'T'};  // Weak
    patternCharMap['S'] = {'C', 'G'};  // Strong
    patternCharMap['B'] = {'C', 'G', 'T'};  // not A
    patternCharMap['D'] = {'A', 'G', 'T'};  // not C
    patternCharMap['H'] = {'A', 'C', 'T'};  // not G
    patternCharMap['V'] = {'A', 'C', 'G'};  // not T
    patternCharMap['N'] = {'A', 'C', 'G', 'T'};  // aNy
}

// Check if pattern character matches sequence character
// Returns false if either character is invalid
bool matchesPattern(char patternChar, char seqChar) {
    // Check if pattern character is valid
    auto it = patternCharMap.find(patternChar);
    if (it == patternCharMap.end()) {
        // Invalid pattern character (not a valid IUPAC code)
        return false;
    }

    // Check if the sequence character is one of the valid DNA bases
    if (seqChar != 'A' && seqChar != 'C' && seqChar != 'G' && seqChar != 'T') {
        // Invalid sequence character
        return false;
    }

    // Check if the pattern allows this sequence character
    return it->second.count(seqChar) > 0;
}

// Validate that sequence contains only valid DNA bases
bool isValidSequence(const string& sequence) {
    for (char c : sequence) {
        if (c != 'A' && c != 'C' && c != 'G' && c != 'T') {
            return false;
        }
    }
    return true;
}

// Validate that pattern contains only valid IUPAC codes
bool isValidPattern(const string& pattern) {
    for (char c : pattern) {
        if (patternCharMap.find(c) == patternCharMap.end()) {
            return false;
        }
    }
    return true;
}

// Check if two pattern characters can potentially match (for LPS building)
bool patternCharsMatch(char p1, char p2) {
    auto it1 = patternCharMap.find(p1);
    auto it2 = patternCharMap.find(p2);

    if (it1 == patternCharMap.end() || it2 == patternCharMap.end()) {
        return false;
    }

    // Check if there's any overlap in possible bases
    for (char base : it1->second) {
        if (it2->second.count(base) > 0) {
            return true;
        }
    }
    return false;
}

// ============================================================================
// NAIVE ALGORITHM - O(n*m) with backtracking
// ============================================================================
bool searchDNASequenceNaive(const string& sequence, const string& pattern) {
    int seqLen = sequence.length();
    int patLen = pattern.length();

    int i = 0;
    int j = 0;

    while (i < seqLen && j < patLen) {
        char seqChar = sequence[i];
        char patChar = pattern[j];

        if (matchesPattern(patChar, seqChar)) {
            i++;
            j++;
        } else {
            i = i - j + 1;  // Backtrack - this is inefficient
            j = 0;
        }
    }
    return j == patLen;
}

// ============================================================================
// KMP ALGORITHM - O(n+m) with preprocessing
// ============================================================================

// Build KMP failure table (LPS array)
vector<int> buildKMPTable(const string& pattern) {
    int patLen = pattern.length();
    vector<int> lps(patLen, 0);
    int length = 0;  // length of previous longest prefix suffix
    int i = 1;

    lps[0] = 0;  // lps[0] is always 0

    while (i < patLen) {
        if (patternCharsMatch(pattern[i], pattern[length])) {
            length++;
            lps[i] = length;
            i++;
        } else {
            if (length != 0) {
                length = lps[length - 1];
            } else {
                lps[i] = 0;
                i++;
            }
        }
    }
    return lps;
}

bool searchDNASequenceKMP(const string& sequence, const string& pattern) {
    int seqLen = sequence.length();
    int patLen = pattern.length();

    if (patLen == 0) {
        return true;
    }
    if (seqLen == 0) {
        return false;
    }

    // Build KMP failure table
    vector<int> lps = buildKMPTable(pattern);

    int i = 0;  // index for sequence
    int j = 0;  // index for pattern

    while (i < seqLen) {
        char seqChar = sequence[i];
        char patChar = pattern[j];

        if (matchesPattern(patChar, seqChar)) {
            i++;
            j++;
        } else {
            if (j != 0) {
                // Use KMP table to avoid redundant comparisons
                j = lps[j - 1];
            } else {
                // No match at all, move to next sequence character
                i++;
            }
        }

        // Found a complete match
        if (j == patLen) {
            return true;
        }
    }

    return false;
}

// ============================================================================
// SLIDING WINDOW ALGORITHM - O(n*m) but simple and cache-friendly
// ============================================================================
bool searchDNASequenceSlidingWindow(const string& sequence, const string& pattern) {
    int seqLen = sequence.length();
    int patLen = pattern.length();

    // Edge cases
    if (patLen == 0) {
        return true;
    }
    if (seqLen < patLen) {
        return false;
    }

    // Slide the window of size patLen across the sequence
    for (int i = 0; i <= seqLen - patLen; i++) {
        bool matched = true;

        // Check if all characters in the current window match the pattern
        for (int j = 0; j < patLen; j++) {
            char seqChar = sequence[i + j];
            char patChar = pattern[j];

            if (!matchesPattern(patChar, seqChar)) {
                matched = false;
                break;  // Early exit on first mismatch
            }
        }

        if (matched) {
            return true;
        }
    }

    return false;
}

// ============================================================================
// MAIN - Test all three algorithms
// ============================================================================
int main() {
    // Initialize the IUPAC pattern map
    initializePatternMap();

    // Test sequences
    vector<string> sequences = {"GATTACA", "GATTG"};
    string pattern = "GATTR";

    cout << "=== Testing All Three Algorithms ===" << endl;
    cout << "Pattern: " << pattern << "\n" << endl;

    // Test with all three algorithms
    for (const string& seq : sequences) {
        bool naiveResult = searchDNASequenceNaive(seq, pattern);
        bool kmpResult = searchDNASequenceKMP(seq, pattern);
        bool slidingWindowResult = searchDNASequenceSlidingWindow(seq, pattern);

        cout << "Sequence: " << seq << endl;
        cout << "  Naive algorithm:          " << (naiveResult ? "true" : "false") << endl;
        cout << "  KMP algorithm:            " << (kmpResult ? "true" : "false") << endl;
        cout << "  Sliding Window algorithm: " << (slidingWindowResult ? "true" : "false") << endl;
        cout << endl;
    }

    // Additional test cases
    cout << "=== Additional Test Cases ===" << endl;

    struct TestCase {
        string seq;
        string pat;
        string desc;
    };

    vector<TestCase> testCases = {
        {"AAAAAAAT", "AAAA", "Pattern with repeats (where KMP shines)"},
        {"AAAAG", "AAAR", "Pattern with wildcards"},
        {"GATTAGA", "GATTNR", "Complex pattern"},
        {"AAATTTGGG", "CCCC", "No match"}
    };

    for (const auto& tc : testCases) {
        bool naiveResult = searchDNASequenceNaive(tc.seq, tc.pat);
        bool kmpResult = searchDNASequenceKMP(tc.seq, tc.pat);
        bool slidingWindowResult = searchDNASequenceSlidingWindow(tc.seq, tc.pat);

        cout << tc.desc << endl;
        cout << "  Seq: " << tc.seq << ", Pat: " << tc.pat << endl;
        cout << "  Naive: " << (naiveResult ? "true" : "false")
             << ", KMP: " << (kmpResult ? "true" : "false")
             << ", Sliding Window: " << (slidingWindowResult ? "true" : "false")
             << "\n" << endl;
    }

    // Test invalid input handling
    cout << "=== Invalid Input Tests ===" << endl;

    // Test 1: Invalid sequence character
    {
        string invalidSeq = "GATTXCA";
        string validPat = "GATT";
        cout << "Invalid sequence (contains X): " << invalidSeq << endl;
        cout << "  Is valid sequence: " << (isValidSequence(invalidSeq) ? "true" : "false") << endl;
        cout << "  Search result: " << (searchDNASequenceKMP(invalidSeq, validPat) ? "true" : "false") << endl;
        cout << endl;
    }

    // Test 2: Invalid pattern character
    {
        string validSeq = "GATTACA";
        string invalidPat = "GATTX";
        cout << "Invalid pattern (contains X): " << invalidPat << endl;
        cout << "  Is valid pattern: " << (isValidPattern(invalidPat) ? "true" : "false") << endl;
        cout << "  Search result: " << (searchDNASequenceKMP(validSeq, invalidPat) ? "true" : "false") << endl;
        cout << endl;
    }

    // Test 3: Valid inputs
    {
        string validSeq = "GATTACA";
        string validPat = "GATTR";
        cout << "Valid inputs - Seq: " << validSeq << ", Pat: " << validPat << endl;
        cout << "  Is valid sequence: " << (isValidSequence(validSeq) ? "true" : "false") << endl;
        cout << "  Is valid pattern: " << (isValidPattern(validPat) ? "true" : "false") << endl;
        cout << "  Search result: " << (searchDNASequenceKMP(validSeq, validPat) ? "true" : "false") << endl;
        cout << endl;
    }

    return 0;
}
