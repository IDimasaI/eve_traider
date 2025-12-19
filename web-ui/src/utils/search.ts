interface SearchResult {
  id: number;
  item: string;
  score: number;
  exactMatch?: boolean;
}

interface SearchConfig {
  minScore?: number;
  exactMatchBonus?: number;
  maxResults?: number; // Ограничение количества результатов
  exactMatchPriority?: boolean; // Приоритет точным совпадениям
}
import { type Items } from './API';
export class FuzzySearcher {

  private preprocessed: Array<{
    id: number;
    original: string;
    lowercase: string;
    words: string[];
  }>;

  constructor(items: Items) {  
    this.preprocessed = items.map((item) => ({
      id: item.id,
      original: item.name,
      lowercase: item.name.toLowerCase(),
      words: item.name.toLowerCase().split(/\s+/),
    }));
  }

  search(query: string, config: SearchConfig = {}): SearchResult[] {
    const {
      minScore = 0.4,
      exactMatchBonus = 0.3,
      maxResults = 10, // По умолчанию показываем топ-10
      exactMatchPriority = true, // Включен приоритет точных совпадений
    } = config;

    const queryWords = query
      .toLowerCase()
      .split(/\s+/)
      .filter((w) => w.length > 0);

    if (queryWords.length === 0) {
      return [];
    }

    // Сначала проверяем точные совпадения
    if (exactMatchPriority) {
      const exactMatches = this.findExactMatches(query, queryWords);
      if (exactMatches.length > 0) {
        return exactMatches.slice(0, maxResults);
      }
    }

    // Если точных совпадений нет, ищем нечёткие
    const results: SearchResult[] = this.preprocessed
      .map((data) => {
        const score = this.calculateMatchScore(queryWords, data);
        const exactMatch = this.isExactMatch(data.original, query);
        return {
          id: data.id,
          item: data.original,
          score: exactMatch ? score + exactMatchBonus : score,
          exactMatch,
        };
      })
      .filter((result) => result.score >= minScore)
      .sort((a, b) => {
        // Сначала сортируем по exactMatch, потом по score
        if (a.exactMatch && !b.exactMatch) return -1;
        if (!a.exactMatch && b.exactMatch) return 1;
        return b.score - a.score;
      });

    // Ограничиваем количество результатов
    return results.slice(0, maxResults);
  }

  private findExactMatches(
    query: string,
    queryWords: string[],
  ): SearchResult[] {
    const exactMatches: SearchResult[] = [];

    for (const data of this.preprocessed) {
      // Проверяем точное совпадение всей строки (без учёта регистра)
      if (data.lowercase === query.toLowerCase()) {
        exactMatches.push({
           id: data.id,
          item: data.original,
          score: 1.0,
          exactMatch: true,
        });
        continue;
      }

      // Проверяем точное совпадение всех слов запроса
      const hasAllWords = queryWords.every((queryWord) =>
        data.words.some((itemWord) => itemWord === queryWord),
      );

      if (hasAllWords) {
        exactMatches.push({
          id: data.id,
          item: data.original,
          score: 0.9,
          exactMatch: true,
        });
      }
    }

    return exactMatches.sort((a, b) => b.score - a.score);
  }

  private isExactMatch(item: string, query: string): boolean {
    const itemLower = item.toLowerCase();
    const queryLower = query.toLowerCase();

    return (
      itemLower === queryLower ||
      itemLower.includes(` ${queryLower} `) ||
      itemLower.startsWith(`${queryLower} `) ||
      itemLower.endsWith(` ${queryLower}`)
    );
  }

  private calculateMatchScore(
    queryWords: string[],
    itemData: { lowercase: string; words: string[] },
  ): number {
    let totalScore = 0;
    let wordsMatched = 0;

    for (const queryWord of queryWords) {
      let bestScore = 0;

      // Проверяем всю строку
      bestScore = Math.max(
        bestScore,
        this.wordInStringScore(queryWord, itemData.lowercase),
      );

      // Проверяем отдельные слова
      for (const itemWord of itemData.words) {
        bestScore = Math.max(
          bestScore,
          this.wordMatchScore(queryWord, itemWord),
        );
      }

      if (bestScore > 0) {
        totalScore += bestScore;
        wordsMatched++;
      }
    }

    // Учитываем процент совпавших слов
    const coverage = wordsMatched / queryWords.length;
    return (totalScore / queryWords.length) * coverage;
  }

  private wordInStringScore(queryWord: string, targetString: string): number {
    if (targetString.includes(queryWord)) return 1.0;

    // Проверка на нечёткое вхождение
    let searchIdx = 0;
    for (
      let i = 0;
      i < targetString.length && searchIdx < queryWord.length;
      i++
    ) {
      if (targetString[i] === queryWord[searchIdx]) {
        searchIdx++;
      }
    }

    return searchIdx / queryWord.length;
  }

  private wordMatchScore(queryWord: string, targetWord: string): number {
    if (targetWord === queryWord) return 1.2;
    if (targetWord.includes(queryWord)) return 1.0;
    if (queryWord.includes(targetWord))
      return targetWord.length / queryWord.length;

    return this.wordInStringScore(queryWord, targetWord);
  }
}
