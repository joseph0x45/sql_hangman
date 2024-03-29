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
  number_of_letters integer not null default 0,
  finished bool not null default false
);

create table guesses (
  id serial primary key,
  guess text not null,
  is_right bool not null,
  game_id integer not null references games(id)
);

CREATE OR REPLACE FUNCTION start_game()
RETURNS TABLE(game_id INTEGER, number_of_letters INTEGER) AS $$
DECLARE 
  new_game_id INTEGER;
  random_word TEXT;
  word_length INTEGER;
BEGIN
  SELECT word INTO random_word FROM words ORDER BY RANDOM() LIMIT 1;
  INSERT INTO games(word_to_guess)
  VALUES(random_word)
  RETURNING id INTO new_game_id;
  SELECT LENGTH(random_word) INTO word_length;
  UPDATE games SET number_of_letters = word_length WHERE id = new_game_id;
  RETURN QUERY SELECT new_game_id, word_length;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_guess(guess TEXT, is_right BOOL, game_id INTEGER)
RETURNS void as $$
BEGIN
  INSERT INTO guesses(guess, is_right, game_id) VALUES(guess, is_right, game_id);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION process_guess(guessed_letter TEXT, game_id INTEGER)
RETURNS INTEGER AS $$
DECLARE 
  current_game record;
  guess_is_right bool;
BEGIN
  SELECT * INTO current_game FROM games where id = game_id;
  IF position(guessed_letter IN current_game.word_to_guess) > 0 THEN
    guess_is_right := true;
  ELSE
    guess_is_right := false;
    PERFORM insert_guess(guessed_letter, guess_is_right, game_id);
  END IF;
  RETURN 2;
END;
$$ LANGUAGE plpgsql;
