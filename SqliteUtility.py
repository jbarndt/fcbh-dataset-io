# SqliteUtility.py
#
# This class provides convenience methods wrapping the Sqlite3 client.
# This class supports the same interface as SQLUtility whenever possible
#

import os
import sys
import sqlite3


class SqliteUtility:

	def getDBPath(database):
		directory = os.environ.get('FCBH_DATASET_DB')
		if directory == None:
			return database
		else:
			return os.path.join(directory, database)


	def destroyDatabase(database):
		databasePath = SqliteUtility.getDBPath(database)
		if os.path.exists(databasePath):
			os.remove(databasePath)


	def __init__(self, database):
		databasePath = SqliteUtility.getDBPath(database)
		self.conn = sqlite3.connect(databasePath)
		self.execute("PRAGMA foreign_keys = ON", ())


	def close(self):
		if self.conn != None:
			self.conn.close()
			self.conn = None


	def execute(self, statement, values=[]):
		cursor = self.conn.cursor()
		try :
			cursor.execute(statement, values)
			self.conn.commit()
			cursor.close()
		except Exception as err:
			self.error(cursor, statement, err)


	def executeInsert(self, statement, values=[]):
		cursor = self.conn.cursor()
		try :
			cursor.execute(statement, values)
			self.conn.commit()
			lastRow = cursor.lastrowid
			cursor.close()
			return lastRow
		except Exception as err:
			self.error(cursor, statement, err)
			return None


	def executeBatch(self, statement, valuesList):
		cursor = self.conn.cursor()
		try:
			cursor.executemany(statement, valuesList)
			self.conn.commit()
			cursor.close()
		except Exception as err:
			self.error(cursor, statement, err)


	def select(self, statement, values=[]):
		cursor = self.conn.cursor()
		try:
			cursor.execute(statement, values)
			resultSet = cursor.fetchall()
			cursor.close()
			return resultSet
		except Exception as err:
			self.error(cursor, statement, err)


	def selectScalar(self, statement, values):
		#print("SQL:", statement, values)
		cursor = self.conn.cursor()
		try:
			cursor.execute(statement, values)
			result = cursor.fetchone()
			cursor.close()
			return result[0] if result != None else None
		except Exception as err:
			self.error(cursor, statement, err)


#	def selectRow(self, statement, values):
#		resultSet = self.select(statement, values)
#		return resultSet[0] if len(resultSet) > 0 else None


#	def selectSet(self, statement, values):
#		resultSet = self.select(statement, values)
#		results = set()
#		for row in resultSet:
#			results.add(row[0])
#		return results		


#	def selectList(self, statement, values):
#		resultSet = self.select(statement, values)
#		results = []
#		for row in resultSet:
#			results.append(row[0])
#		return results


#	def selectMap(self, statement, values):
#		resultSet = self.select(statement, values)
#		results = {}
#		for row in resultSet:
#			results[row[0]] = row[1]
#		return results


#	def selectMapList(self, statement, values):
#		resultSet = self.select(statement, values)
#		results = {}
#		for row in resultSet:
#			values = results.get(row[0], [])
#			values.append(row[1])
#			results[row[0]] = values
#		return results


	def error(self, cursor, stmt, error):
		cursor.close()	
		print("ERROR executing SQL %s on '%s'" % (error, stmt))
		self.conn.rollback()
		sys.exit()


if __name__ == "__main__":
	database = "TestSqliteUtility.db"
	SqliteUtility.destroyDatabase(database)
	sql = SqliteUtility(database)
	sql.execute("CREATE TABLE abc (a int, b real, c text)", [])
	sql.execute("INSERT INTO abc (a, b, c) VALUES (1, 1.0, '1')", [])
	sql.execute("INSERT INTO abc (a, b, c) VALUES (?,?,?)", [2, 2.2, '2'])
	records = []
	records.append([3, 3.3, '3'])
	records.append([4, 4.4, '4'])
	sql.executeBatch("INSERT INTO abc(a, b, c) VALUES (?,?,?)", records)
	results = sql.select("SELECT a, b, c FROM abc WHERE a < ?", [10])
	print(results)
	sql.close(), 
