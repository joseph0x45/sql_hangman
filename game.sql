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
BEGIN
  SELECT word
  INTO random_word
  FROM words ORDER BY RANDOM() LIMIT 1;
  INSERT INTO games(word_to_guess)
  VALUES(random_word)
  RETURNING id INTO new_game_id;
  UPDATE games SET number_of_letters = LENGTH(random_word) WHERE id = new_game_id;
  RETURN QUERY SELECT new_game_id, LENGTH(random_word);
END;
$$ LANGUAGE plpgsql;
