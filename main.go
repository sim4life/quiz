package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type ByLength []string

func (s ByLength) Len() int           { return len(s) }
func (s ByLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByLength) Less(i, j int) bool { return len(s[i]) < len(s[j]) }

type StringSet struct {
	set map[string]bool
}

func NewStringSet() *StringSet {
	return &StringSet{make(map[string]bool)}
}

func (set *StringSet) Add(s string) bool {
	_, found := set.set[s]
	set.set[s] = true
	return !found //False if it existed already
}

func (set *StringSet) Get(s string) bool {
	_, found := set.set[s]
	return found //true if it existed already
}

func (set *StringSet) Remove(s string) {
	delete(set.set, s)
}

func (set *StringSet) isEmpty() bool {
	for _, val := range set.set {
		if val == true {
			return false
		}
	}
	return true
}

func isSearchableWord(word string, min_len int) bool {
	return word != "" && len(word) >= min_len
}

func isPrefixFound(word string, found bool) bool {
	return word == "" || (word != "" && found)
}

/*
 * The getCompoundingWords collects all the component words of the search_word
 */
func getCompoundingWords(sorted_words []string, search_word string) ([]string, int) {
	comp_words := make([]string, 0)
	min_word_len := len(sorted_words[0])
	for j := len(sorted_words) - 1; j >= 0; j-- {
		//assuming no duplicates and starting search from the index of the word with
		//length of the search_word minus min word length
		if strings.Contains(search_word, sorted_words[j]) {
			if len(search_word) >= (len(sorted_words[j]) + min_word_len) {
				comp_words = append(comp_words, sorted_words[j])
			}
		}
	}
	return comp_words, min_word_len
}

/*
 * The searchSubWord attempts to match search_word in the lex_sort_words list.
 * If it can't search it then it searches for prefix, middle words and suffix
 * of search_word by calling function isStringConcat
 */
func searchSubWord(sorted_words, lex_sort_words []string, search_word string,
	startIndices []int, foundSet *StringSet, total_words, min_word_len int, found bool) bool {
	pos := sort.SearchStrings(lex_sort_words, search_word)
	if pos < total_words && search_word == lex_sort_words[pos] {
		found = true
	} else {
		if len(search_word) > min_word_len {
			if foundSet.Get(search_word) {
				found = true
			} else {
				found = isStringConcat(sorted_words[:startIndices[len(search_word)-min_word_len]+1],
					lex_sort_words, search_word, startIndices, foundSet)
				if found {
					foundSet.Add(search_word)
				}
			}
		}
	}

	return found
}

/**
 * The isStringConcat is a recursive function and it return true if
 * it finds all the components of the search_word in the sorted_words list.
 */
func isStringConcat(sorted_words, lex_sort_words []string, search_word string,
	startIndices []int, foundSet *StringSet) bool {
	total_words := len(lex_sort_words)

	comp_words, min_word_len := getCompoundingWords(sorted_words, search_word)

	//Find full component words within the list
	for _, c_word := range comp_words {
		found := false

		tokens := strings.Split(search_word, c_word)
		// if prefix is NOT a searchable word then c_word is a prefix
		if isSearchableWord(tokens[0], min_word_len) {
			// if non-Prefix is a searchable word then we need to run search on Prefix
			if isSearchableWord(tokens[1], min_word_len) {
				found = searchSubWord(sorted_words, lex_sort_words, tokens[0], startIndices,
					foundSet, total_words, min_word_len, found)
			}
		}
		/* if non-Prefix is NOT searchable word then c_word is NOT a suffix
		and/or it's a middle word */
		if isSearchableWord(tokens[1], min_word_len) {
			// run a search on non-Prefix word only if prefix is found
			if isPrefixFound(tokens[0], found) {
				found = searchSubWord(sorted_words, lex_sort_words, tokens[1], startIndices,
					foundSet, total_words, min_word_len, found)
			}
		}
		if found {
			return true
		}
	}

	return false
}

/*
 * The getLongestCompWord is a non-recursive function.
 * It prepares and optimizes parameters to be used by the recursive function.
 * Then it calls the recursive function on each word in the word list
 * starting from the longest word in the word list
 */
func getLongestCompWord(words []string) (longest_comp_word string) {

	longest_comp_word = ""
	found := false
	total_words := len(words)

	foundSet := NewStringSet()

	//preparing words list in ascending order lexical-wise
	lex_sort_words := make([]string, len(words))
	copy(lex_sort_words, words)
	sort.Strings(lex_sort_words)
	//preparing words list in ascending order length-wise
	sorted_words := make([]string, len(words))
	copy(sorted_words, words)
	sort.Sort(ByLength(sorted_words))

	//slice of startIndices of each diff len of word in ascending order words list
	//for optimization purposes in the recursive function
	startIndices := make([]int, len(sorted_words[0])+1)

	for i := 1; i < total_words; i++ {
		if len(sorted_words[i]) > len(sorted_words[i-1]) {
			for j := len(sorted_words[i-1]) + 1; j < len(sorted_words[i]); j++ {
				//fill the gap for non-contagious word length increase
				startIndices = append(startIndices, startIndices[len(sorted_words[i-1])])
			}
			startIndices = append(startIndices, i)
		}
	}

	fmt.Printf("[%s] Evaluating word[%d], sorted length-wise is: %s ...\n",
		time.Now().Format(time.Stamp), total_words-1, sorted_words[total_words-1])
	//sending each word from the sorted word list length-wise for evaluation
	//starting from the biggest word
	for i := total_words - 1; i >= 0; i-- {
		//fmt.Printf("[%s] Evaluating word[%d], sorted length-wise is: %s ...\n", time.Now().Format(time.Stamp), i, sorted_words[i])
		found = isStringConcat(sorted_words, lex_sort_words, sorted_words[i], startIndices, foundSet)
		if found {
			fmt.Printf("[%s] Evaluating & FOUND word[%d], sorted length-wise is: %s ...\n",
				time.Now().Format(time.Stamp), i, sorted_words[i])
			longest_comp_word = sorted_words[i]
			fmt.Println("found word is:", sorted_words[i])
			return
		}
	}
	return
}

/*
 * The main function manages input parameters,
 * reads word list file into word list slice and
 * calls the getLongestComWord non-recursive function
 * to get the longest compound word.
 */
func main() {
	file_path, err := checkArgs()
	if err != "" {
		fmt.Printf("Error: %s", err)

		file_path = "word.list"
		fmt.Printf("No input file found...going to read: %s\n\n", file_path)

	}

	lines, err2 := readLines(file_path)
	if err2 != nil {
		log.Fatalf("readLines: %s", err2)
	}

	longest_comp_word := getLongestCompWord(lines)
	if longest_comp_word == "" {
		fmt.Println("No compound word found in the list")
	} else {
		fmt.Println("The longest compound word in the list is: ", longest_comp_word)
	}
}

/*
 * The readLines function reads a whole file into memory
 * and returns a slice of its lines.
 */
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

/*
 * The checkArgs() function returns a string of file path and
 * a string of error if there is any.
 */
func checkArgs() (string, string) {
	//Fetch the command line arguments.
	args := os.Args

	//Check the length of the arugments, return failure if that are too
	//long or too short.
	if (len(args) < 2) || (len(args) >= 3) {
		return "1", "Invalid number of arguments. \n" +
			"Please provide the file name with relative path of the words list input file!\n"
	}
	file_path := args[1]
	//On success, return the file_path value and an empty string indicating
	//that everything is good.
	return file_path, ""
}
