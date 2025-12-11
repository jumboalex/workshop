#include <iostream>
#include <string>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <algorithm>

using namespace std;

// IUPAC nucleotide codes mapping
unordered_map<char, unordered_set<char>> patternCharMap;

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

// DNASearchEngine handles multiple DNA sequences with n-gram indexing
class DNASearchEngine {
private:
    vector<string> sequences;
    unordered_map<string, unordered_set<string>> ngramIndex; // ngram -> set of sequences
    int n; // n-gram size

    // Extract all n-grams from a sequence
    vector<string> extractNGrams(const string& sequence) const {
        vector<string> ngrams;
        if (sequence.length() < n) {
            return ngrams;
        }

        for (size_t i = 0; i <= sequence.length() - n; i++) {
            ngrams.push_back(sequence.substr(i, n));
        }
        return ngrams;
    }

    // Check if sequence contains query (with wildcard support)
    bool matchesQuery(const string& sequence, const string& query) const {
        if (query.length() > sequence.length()) {
            return false;
        }

        // Sliding window search with wildcard matching
        for (size_t i = 0; i <= sequence.length() - query.length(); i++) {
            bool matched = true;
            for (size_t j = 0; j < query.length(); j++) {
                char seqChar = sequence[i + j];
                char queryChar = query[j];

                // Check if query character matches sequence character (with wildcards)
                auto it = patternCharMap.find(queryChar);
                if (it == patternCharMap.end() || it->second.count(seqChar) == 0) {
                    matched = false;
                    break;
                }
            }
            if (matched) {
                return true;
            }
        }
        return false;
    }

    // Direct search for short queries
    vector<string> directSearch(const string& query) const {
        vector<string> results;
        for (const auto& seq : sequences) {
            if (matchesQuery(seq, query)) {
                results.push_back(seq);
            }
        }
        return results;
    }

public:
    DNASearchEngine(const vector<string>& seqs, int ngramSize)
        : sequences(seqs), n(ngramSize) {
        buildIndex();
    }

    // Build n-gram index from all sequences
    void buildIndex() {
        cout << "=== Building N-Gram Index ===" << endl;
        cout << "N-gram size: " << n << endl;
        cout << "Sequences to index: " << sequences.size() << "\n" << endl;

        for (const auto& seq : sequences) {
            vector<string> ngrams = extractNGrams(seq);

            cout << "Sequence: " << seq << endl;
            cout << "  N-grams: [";
            for (size_t i = 0; i < ngrams.size(); i++) {
                if (i > 0) cout << ", ";
                cout << ngrams[i];
            }
            cout << "]" << endl;

            for (const auto& ngram : ngrams) {
                ngramIndex[ngram].insert(seq);
            }
        }

        cout << "\n=== N-Gram Index ===" << endl;
        for (const auto& entry : ngramIndex) {
            cout << entry.first << " -> [";
            bool first = true;
            for (const auto& seq : entry.second) {
                if (!first) cout << ", ";
                cout << seq;
                first = false;
            }
            cout << "]" << endl;
        }
        cout << endl;
    }

    // Search for sequences matching query
    vector<string> search(const string& query) {
        cout << "\n=== Searching for Query: " << query << " ===" << endl;

        // Step 1: Extract n-grams from query
        vector<string> queryNGrams = extractNGrams(query);

        cout << "Query n-grams: [";
        for (size_t i = 0; i < queryNGrams.size(); i++) {
            if (i > 0) cout << ", ";
            cout << queryNGrams[i];
        }
        cout << "]\n" << endl;

        if (queryNGrams.empty()) {
            // Query is shorter than n-gram size, fall back to direct search
            return directSearch(query);
        }

        // Step 2: Look up each n-gram and collect candidate sequences
        cout << "N-gram lookups:" << endl;
        unordered_set<string> candidates;

        for (size_t i = 0; i < queryNGrams.size(); i++) {
            const string& ngram = queryNGrams[i];
            auto it = ngramIndex.find(ngram);

            cout << "  " << ngram << " -> [";
            if (it != ngramIndex.end()) {
                bool first = true;
                for (const auto& seq : it->second) {
                    if (!first) cout << ", ";
                    cout << seq;
                    first = false;
                }
            }
            cout << "]" << endl;

            if (i == 0) {
                // Initialize with first n-gram's sequences
                if (it != ngramIndex.end()) {
                    candidates = it->second;
                }
            } else {
                // Intersect with current candidates
                if (it != ngramIndex.end()) {
                    unordered_set<string> newCandidates;
                    for (const auto& seq : candidates) {
                        if (it->second.count(seq) > 0) {
                            newCandidates.insert(seq);
                        }
                    }
                    candidates = newCandidates;
                } else {
                    candidates.clear();
                }
            }
        }

        // Step 3: Get candidate list
        cout << "\nCandidates after intersection: [";
        bool first = true;
        for (const auto& cand : candidates) {
            if (!first) cout << ", ";
            cout << cand;
            first = false;
        }
        cout << "]" << endl;

        // Step 4: Filter false positives - verify full match
        cout << "\nVerifying candidates:" << endl;
        vector<string> results;
        for (const auto& seq : candidates) {
            if (matchesQuery(seq, query)) {
                cout << "  " << seq << ": ✓ MATCH" << endl;
                results.push_back(seq);
            } else {
                cout << "  " << seq << ": ✗ FALSE POSITIVE" << endl;
            }
        }

        return results;
    }
};

int main() {
    initializePatternMap();

    // Upload DNA sequences
    vector<string> sequences = {"GATTACA", "GATTG"};

    // Create search engine with n-gram size 4
    DNASearchEngine engine(sequences, 4);

    // Test cases
    struct TestCase {
        string query;
        string description;
    };

    vector<TestCase> testCases = {
        {"GATT", "Exact match at beginning"},
        {"ATTACA", "Exact match in middle/end"},
        {"GATTR", "Wildcard R (A or G)"},
        {"GATTM", "Wildcard M (A or C)"},
        {"GATTRR", "Double wildcard RR (no match expected)"}
    };

    cout << "\n" << string(70, '=') << endl;
    cout << "RUNNING SEARCH TESTS" << endl;
    cout << string(70, '=') << endl;

    for (const auto& tc : testCases) {
        cout << "\n" << string(70, '-') << endl;
        cout << "Test: " << tc.description << endl;
        vector<string> results = engine.search(tc.query);

        cout << "\nRESULT: Matching sequences = [";
        for (size_t i = 0; i < results.size(); i++) {
            if (i > 0) cout << ", ";
            cout << results[i];
        }
        cout << "]" << endl;
    }

    // Demonstrate false positive example
    cout << "\n" << string(70, '=') << endl;
    cout << "FALSE POSITIVE EXAMPLE" << endl;
    cout << string(70, '=') << endl;

    vector<string> falsePositiveSeqs = {"ATTAGATT"};
    DNASearchEngine fpEngine(falsePositiveSeqs, 4);
    vector<string> fpResults = fpEngine.search("GATTA");

    cout << "\nRESULT: Matching sequences = [";
    for (size_t i = 0; i < fpResults.size(); i++) {
        if (i > 0) cout << ", ";
        cout << fpResults[i];
    }
    cout << "]" << endl;
    cout << "(Should be empty - GATTA is not in ATTAGATT as contiguous substring)" << endl;

    return 0;
}
