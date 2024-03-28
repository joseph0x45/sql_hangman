create table words (
  word text not null
);

INSERT INTO words (word) VALUES
('apple'),
('banana'),
('cat'),
('dog'),
('elephant'),
('fish'),
('grape'),
('house'),
('xenophobia'),
('jacket');

create table games (
  id serial primary key,
  word_to_guess text not null,
  number_of_letters integer not null,
  finished bool not null default false
);

create table guesses (
  id serial primary key,
  guess text not null,
  is_right bool not null,
  game_id integer not null references games(id)
);

CREATE OR REPLACE FUNCTION start_game()
RETURNS TRIGGER AS $$
BEGIN
  WITH random_word_data AS (
    SELECT word, LENGTH(word) AS word_length
    FROM words ORDER BY RANDOM() LIMIT 1
  )
  INSERT INTO games(word_to_guess, number_of_letters)
  SELECT word, word_length FROM random_word_data;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_guess(guess_id INTEGER, guessed_letter TEXT, game_id INTEGER)
RETURNS TRIGGER AS $$
BEGIN
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_Guesses_AfterInsert
AFTER INSERT ON guesses
FOR EACH ROW
EXECUTE FUNCTION check_guess(NEW.id, NEW.guess, NEW.game_id);
